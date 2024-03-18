package main

import (
	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/utils"
	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/WangYihang/http-grab/pkg/option"
)

func main() {
	scheduler := gojob.NewScheduler().
		SetNumWorkers(option.Opt.NumWorkers).
		SetMaxRetries(option.Opt.MaxTries).
		SetMaxRuntimePerTaskSeconds(option.Opt.MaxRuntimePerTaskSeconds).
		SetNumShards(int64(option.Opt.NumShards)).
		SetShard(int64(option.Opt.Shard)).
		SetOutputFilePath(option.Opt.OutputFilePath).
		Start()
	for line := range utils.Cat(option.Opt.InputFilePath) {
		scheduler.Submit(model.NewTask(line, option.Opt.Port, option.Opt.Path, option.Opt.Host))
	}
	scheduler.Wait()
}
