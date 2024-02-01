package main

import (
	"log/slog"
	"os"

	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFilePath         string `short:"i" long:"input" description:"input file path" required:"true"`
	OutputFilePath        string `short:"o" long:"output" description:"output file path" required:"true"`
	StatusUpdatesFilePath string `short:"s" long:"status-updates" description:"status updates file path"`

	NumWorkers int   `short:"n" long:"num-workers" description:"number of workers" default:"32"`
	NumShards  int64 `long:"num-shards" description:"number of shards" default:"1"`
	Shard      int64 `long:"shard" description:"shard" default:"0"`

	Port     int    `short:"p" long:"port" description:"port" default:"80"`
	Path     string `long:"path" description:"path" default:"index.html"`
	Host     string `long:"host" description:"http host header" default:""`
	MaxTries int    `short:"m" long:"max-tries" description:"max tries" default:"4"`
	Timeout  int    `short:"t" long:"timeout" description:"timeout" default:"8"`
}

var opts Options

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	if opts.StatusUpdatesFilePath == "" {
		opts.StatusUpdatesFilePath = opts.OutputFilePath + ".status"
	}
	if opts.Shard >= opts.NumShards {
		slog.Error("shard must be less than num-shards", slog.Int64("shard", opts.Shard), slog.Int64("num_shards", opts.NumShards))
		os.Exit(1)
	}
}

func load() chan *model.Task {
	return model.LoadTasks(
		model.NewTask,
		opts.InputFilePath,
		opts.Port,
		opts.Path,
		opts.Host,
		opts.Timeout,
		opts.MaxTries,
		opts.NumShards,
		opts.Shard,
	)
}

func main() {
	numWorkers := opts.NumWorkers
	numTasks := model.CountLines(opts.InputFilePath, opts.NumShards, opts.Shard)
	slog.Info("all tasks loaded", slog.Int64("num_tasks", numTasks))
	resultChans := make([]chan *model.Task, 0, numWorkers)
	for _, taskChan := range model.FanOut(load(), numWorkers) {
		resultChans = append(resultChans, model.Worker(taskChan))
	}
	model.StoreTasks(model.FanIn(resultChans), opts.OutputFilePath, opts.StatusUpdatesFilePath, numTasks)
}
