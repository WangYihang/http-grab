package main

import (
	"os"
	"time"

	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/runner"
	"github.com/WangYihang/gojob/pkg/utils"
	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/WangYihang/http-grab/pkg/option"
	"github.com/jessevdk/go-flags"
)

var Opt option.Option

func init() {
	Opt.Version = model.PrintVersion
	if _, err := flags.Parse(&Opt); err != nil {
		os.Exit(1)
	}
}

func main() {
	scheduler := gojob.New(
		gojob.WithNumWorkers(Opt.NumWorkers),
		gojob.WithMaxRetries(Opt.MaxTries),
		gojob.WithMaxRuntimePerTaskSeconds(Opt.MaxRuntimePerTaskSeconds),
		gojob.WithNumShards(int64(Opt.NumShards)),
		gojob.WithShard(int64(Opt.Shard)),
		gojob.WithResultFilePath(Opt.OutputFilePath),
		gojob.WithMetadata("build", map[string]string{
			"version": model.Version,
			"commit":  model.Commit,
			"date":    model.Date,
		}),
		gojob.WithMetadata("runner", runner.Runner),
		gojob.WithMetadata("arguments", Opt),
		gojob.WithMetadata("started_at", time.Now().Format(time.RFC3339)),
	).
		Start()
	for line := range utils.Cat(Opt.InputFilePath) {
		scheduler.Submit(
			model.NewTask(line).
				WithPort(Opt.Port).
				WithPath(Opt.Path).
				WithHost(Opt.Host).
				WithSNI(Opt.Host).
				WithScheme(Opt.Scheme).
				WithTimeout(Opt.Timeout),
		)
	}
	scheduler.Wait()
}
