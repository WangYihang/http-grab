package main

import (
	"os"

	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/util"
	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFilePath  string `short:"i" long:"input" description:"input file path" required:"true"`
	OutputFilePath string `short:"o" long:"output" description:"output file path" required:"true"`
	StatusFilePath string `short:"s" long:"status" description:"status file path" required:"true" default:"-"`

	NumWorkers               int   `short:"n" long:"num-workers" description:"number of workers" default:"32"`
	NumShards                int64 `long:"num-shards" description:"number of shards" default:"1"`
	Shard                    int64 `long:"shard" description:"shard" default:"0"`
	MaxTries                 int   `short:"m" long:"max-tries" description:"max tries" default:"4"`
	MaxRuntimePerTaskSeconds int   `short:"t" long:"max-runtime-per-task-seconds" description:"max runtime per task seconds" default:"60"`

	Port int    `short:"p" long:"port" description:"port" default:"80"`
	Path string `long:"path" description:"path" default:"index.html"`
	Host string `long:"host" description:"http host header, leave it blank to use the IP address" default:""`
}

var opts Options

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	scheduler := gojob.NewScheduler().
		SetNumWorkers(opts.NumWorkers).
		SetMaxRetries(opts.MaxTries).
		SetMaxRuntimePerTaskSeconds(opts.MaxRuntimePerTaskSeconds).
		SetNumShards(int64(opts.NumShards)).
		SetShard(int64(opts.Shard)).
		SetOutputFilePath(opts.OutputFilePath).
		Start()
	for line := range util.Cat(opts.InputFilePath) {
		scheduler.Submit(model.NewTask(line, opts.Port, opts.Path, opts.Host))
	}
	scheduler.Wait()
}
