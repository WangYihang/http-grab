package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/runner"
	"github.com/WangYihang/gojob/pkg/utils"
	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/WangYihang/http-grab/pkg/option"
	"github.com/google/uuid"
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
		gojob.WithStatusFilePath(Opt.StatusFilePath),
		gojob.WithMetadataFilePath(Opt.MetadataFilePath),
		gojob.WithTotalTasks(utils.Count(utils.Cat(Opt.InputFilePath))),
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
		parameters := make(map[string]string)
		parameters["ingress"] = strings.TrimSpace(line)
		parameters["timestamp"] = fmt.Sprint(time.Now().UnixMilli())
		parameters["challenge"] = uuid.New().String()
		scheduler.Submit(
			model.NewTask(line).
				WithPort(Opt.Port).
				WithPath(Opt.Path).
				WithHost(Opt.Host).
				WithSNI(Opt.Host).
				WithMethod(Opt.Method).
				WithScheme(Opt.Scheme).
				WithParameters(parameters).
				WithTimeout(Opt.Timeout),
		)
	}
	scheduler.Wait()
}
