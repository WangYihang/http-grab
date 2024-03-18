package http_grab_task

import (
	"fmt"
	"strings"
)

type HTTPGrabArguments struct {
	InputFilePath  string
	OutputFilePath string
	StatusFilePath string

	NumWorkers               int
	NumShards                int
	Shard                    int
	MaxTries                 int
	MaxRuntimePerTaskSeconds int

	Port    int
	Path    string
	Host    string
	Timeout int
}

func NewHTTPGrabArguments() *HTTPGrabArguments {
	return &HTTPGrabArguments{
		InputFilePath:            "input.txt",
		OutputFilePath:           "output.txt",
		StatusFilePath:           "output.status",
		NumWorkers:               1024,
		NumShards:                1,
		Shard:                    0,
		MaxTries:                 2,
		MaxRuntimePerTaskSeconds: 4,
		Port:                     80,
		Path:                     "/",
		Host:                     "localhost",
		Timeout:                  4,
	}
}

func (z *HTTPGrabArguments) WithInputFilePath(inputFilePath string) *HTTPGrabArguments {
	z.InputFilePath = inputFilePath
	return z
}

func (z *HTTPGrabArguments) WithOutputFilePath(outputFilePath string) *HTTPGrabArguments {
	z.OutputFilePath = outputFilePath
	return z
}

func (z *HTTPGrabArguments) WithStatusFilePath(statusFilePath string) *HTTPGrabArguments {
	z.StatusFilePath = statusFilePath
	return z
}

func (z *HTTPGrabArguments) WithNumWorkers(numWorkers int) *HTTPGrabArguments {
	z.NumWorkers = numWorkers
	return z
}

func (z *HTTPGrabArguments) WithNumShards(numShards int) *HTTPGrabArguments {
	z.NumShards = numShards
	return z
}

func (z *HTTPGrabArguments) WithShard(shard int) *HTTPGrabArguments {
	z.Shard = shard
	return z
}

func (z *HTTPGrabArguments) WithMaxTries(maxTries int) *HTTPGrabArguments {
	z.MaxTries = maxTries
	return z
}

func (z *HTTPGrabArguments) WithMaxRuntimePerTaskSeconds(maxRuntimePerTaskSeconds int) *HTTPGrabArguments {
	z.MaxRuntimePerTaskSeconds = maxRuntimePerTaskSeconds
	return z
}

func (z *HTTPGrabArguments) WithPort(port int) *HTTPGrabArguments {
	z.Port = port
	return z
}

func (z *HTTPGrabArguments) WithPath(path string) *HTTPGrabArguments {
	z.Path = path
	return z
}

func (z *HTTPGrabArguments) WithHost(host string) *HTTPGrabArguments {
	z.Host = host
	return z
}

func (z *HTTPGrabArguments) WithTimeout(timeout int) *HTTPGrabArguments {
	z.Timeout = timeout
	return z
}

func (z *HTTPGrabArguments) String() string {
	arguments := []string{
		"--input", z.InputFilePath,
		"--output", z.OutputFilePath,
		"--status", z.StatusFilePath,
		"--num-workers", fmt.Sprintf("%d", z.NumWorkers),
		"--num-shards", fmt.Sprintf("%d", z.NumShards),
		"--shard", fmt.Sprintf("%d", z.Shard),
		"--max-tries", fmt.Sprintf("%d", z.MaxTries),
		"--max-runtime-per-task-seconds", fmt.Sprintf("%d", z.MaxRuntimePerTaskSeconds),
		"--port", fmt.Sprintf("%d", z.Port),
		"--path", z.Path,
		"--host", z.Host,
		"--timeout", fmt.Sprintf("%d", z.Timeout),
	}
	return strings.Join(arguments, " ")
}