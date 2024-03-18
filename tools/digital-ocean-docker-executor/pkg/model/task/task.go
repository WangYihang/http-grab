package http_grab_task

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/executor/secureshell"
	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/task"
	"github.com/WangYihang/http-grab/tools/digital-ocean-docker-executor/pkg/option"
	"github.com/charmbracelet/log"
)

type HTTPGrabTask struct {
	e           *secureshell.SSHExecutor
	image       string
	containerID string
	labels      map[string]interface{}
	arguments   *HTTPGrabArguments
	folder      string
	port        int
}

func Generate(name string, port int) <-chan *HTTPGrabTask {
	out := make(chan *HTTPGrabTask)
	go func() {
		defer close(out)
		shards := 8
		for shard := range shards {
			out <- New(port, shard, shards, name)
		}
	}()
	return out
}

func New(port, shard, shards int, label string) *HTTPGrabTask {
	folder := fmt.Sprintf("/data/%s/shards-%d/shard-%d", label, shards, shard)
	filename := fmt.Sprintf("%s-%d-%d", label, shard, shards)
	z := &HTTPGrabTask{
		arguments: NewHTTPGrabArguments().
			WithInputFilePath("/data/input.txt").
			WithOutputFilePath(fmt.Sprintf("/data/%s.json", filename)).
			WithStatusFilePath(fmt.Sprintf("/data/%s.status", filename)).
			WithTimeout(8).
			WithMaxTries(4).
			WithMaxRuntimePerTaskSeconds(32).
			WithShard(shard).
			WithNumShards(shards).
			WithPort(port).
			WithHost("localhost").
			WithPath("/").
			WithNumWorkers(1024),
		labels: make(map[string]interface{}),
		image:  "ghcr.io/wangyihang/http-grab:main",
		folder: folder,
		port:   port,
	}
	z.labels["task.label"] = label
	z.labels["task.shard"] = z.arguments.Shard
	z.labels["task.num-shards"] = z.arguments.NumShards
	return z
}

func (h *HTTPGrabTask) WithArguments(arguments *HTTPGrabArguments) *HTTPGrabTask {
	h.arguments = arguments
	return h
}

func (h *HTTPGrabTask) String() string {
	return h.arguments.String()
}

func (h *HTTPGrabTask) Assign(e *secureshell.SSHExecutor) error {
	h.e = e
	return h.e.Connect()
}

func (h *HTTPGrabTask) Prepare() error {
	images := []string{
		"amazon/aws-cli",
		h.image,
	}
	for _, image := range images {
		h.e.RunCommand(strings.Join([]string{
			"docker", "pull", image,
		}, " "))
	}
	h.e.RunCommand(fmt.Sprintf(
		"mkdir -p %s && wget -O %s/ipinfo-%d-%d.json https://ipinfo.io/json",
		h.folder,
		h.folder,
		h.arguments.Shard,
		h.arguments.NumShards,
	))
	// upload input file
	h.e.UploadFile("zmap.txt", filepath.Join(h.folder, "input.txt"))
	return nil
}

func (h *HTTPGrabTask) Start() error {
	arguments := []string{
		"docker", "run",
		"--interactive", "--tty", "--detach",
		"--network", "host",
		"--volume", fmt.Sprintf(
			"%s:/data",
			h.folder,
		),
	}
	for k, v := range h.labels {
		arguments = append(arguments, "--label", fmt.Sprintf("%s=%v", k, v))
	}
	arguments = append(arguments, h.image)
	arguments = append(arguments, h.arguments.String())
	stdout, stderr, err := h.e.RunCommand(strings.Join(arguments, " "))
	if err != nil {
		return err
	}
	if stderr != "" {
		return fmt.Errorf(stderr)
	}
	h.containerID = strings.TrimSpace(stdout)
	return nil
}

func (h *HTTPGrabTask) Stop() error {
	_, _, err := h.e.RunCommand(strings.Join([]string{
		h.containerID,
	}, " "))
	return err
}

func (h *HTTPGrabTask) Status() (task.StatusInterface, error) {
	// check if task is running
	arguments := []string{
		"docker", "ps",
		"--all",
		"--quiet",
	}
	for k, v := range h.labels {
		arguments = append(arguments, "--filter", fmt.Sprintf("label=%s=%v", k, v))
	}
	stdout, stderr, err := h.e.RunCommand(strings.Join(arguments, " "))
	log.Info("task status", "stdout", stdout, "stderr", stderr, "err", err, "task", h.String())
	if err != nil {
		return nil, err
	}
	if stderr != "" {
		return nil, fmt.Errorf(stderr)
	}
	if stdout == "" {
		return PendingProgress, nil
	}

	log.Info("task have already been started", "container", stdout)
	h.containerID = strings.TrimSpace(stdout)

	// check if the container is running
	stdout, stderr, err = h.e.RunCommand(strings.Join([]string{
		"docker",
		"inspect",
		"--format",
		"{{.State.Running}}",
		h.containerID,
	}, " "))
	if err != nil {
		return nil, err
	}
	if stderr != "" {
		return nil, fmt.Errorf(stderr)
	}
	if strings.TrimSpace(stdout) == "false" {
		// the container is not running, so it must have finished
		return DoneProgress, nil
	}

	// read the status file
	stdout, stderr, err = h.e.RunCommand(strings.Join([]string{
		"tail",
		"-n",
		"1",
		filepath.Join(h.folder, filepath.Base(h.arguments.StatusFilePath)),
	}, " "))
	if err != nil {
		return PendingProgress, nil
	}
	if stderr != "" {
		return PendingProgress, nil
	}

	// parse the status file
	progress, err := NewHTTPGrabProgress(stdout)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (h *HTTPGrabTask) Download() error {
	// upload to amazon s3
	if option.Opt.S3AccessKey != "" {
		today := time.Now().Format("2006-01-02")
		h.e.RunCommand(fmt.Sprintf("docker run --rm -v ~/.aws:/root/.aws amazon/aws-cli configure set aws_access_key_id %s", option.Opt.S3Option.S3AccessKey))
		h.e.RunCommand(fmt.Sprintf("docker run --rm -v ~/.aws:/root/.aws amazon/aws-cli configure set aws_secret_access_key %s", option.Opt.S3Option.S3SecretKey))
		h.e.RunCommand(fmt.Sprintf("docker run --rm -v ~/.aws:/root/.aws amazon/aws-cli configure set default.region %s", option.Opt.S3Option.S3Region))
		h.e.RunCommand(fmt.Sprintf("docker run --rm -v %s:/data -v ~/.aws:/root/.aws amazon/aws-cli s3 cp /data/ s3://%s/%s/%d/%s/ --recursive", h.folder, option.Opt.S3Bucket, option.Opt.Name, h.port, today))
	}

	// Download to local
	h.e.DownloadFile(filepath.Join("/data", h.arguments.OutputFilePath), filepath.Join("data", filepath.Base(h.arguments.OutputFilePath)))
	h.e.DownloadFile(filepath.Join("/data", h.arguments.StatusFilePath), filepath.Join("data", filepath.Base(h.arguments.StatusFilePath)))

	// Remove files
	h.e.RunCommand(fmt.Sprintf("rm -rf %s", h.folder))
	return nil
}
