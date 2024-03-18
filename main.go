package main

import (
	"os"

	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/utils"
	"github.com/WangYihang/http-grab/pkg/loaders"
	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/WangYihang/http-grab/pkg/option"
	"github.com/jessevdk/go-flags"
)

var Opt option.Option

func init() {
	_, err := flags.Parse(&Opt)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	scheduler := gojob.NewScheduler().
		SetNumWorkers(Opt.NumWorkers).
		SetMaxRetries(Opt.MaxTries).
		SetMaxRuntimePerTaskSeconds(Opt.MaxRuntimePerTaskSeconds).
		SetNumShards(int64(Opt.NumShards)).
		SetShard(int64(Opt.Shard)).
		SetTotalTasks(utils.Count(loaders.Get(Opt.InputFilePath, "txt"))).
		SetOutputFilePath(Opt.OutputFilePath).
		Start()
	for line := range utils.Cat(Opt.InputFilePath) {
		scheduler.Submit(model.NewTask(line, Opt.Port, Opt.Path, Opt.Host, Opt.Timeout))
	}
	scheduler.Wait()
}
