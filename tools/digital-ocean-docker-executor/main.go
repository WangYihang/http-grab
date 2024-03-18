package main

import (
	"os"

	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/provider"
	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/provider/api"
	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/scheduler"
	gojob_utils "github.com/WangYihang/gojob/pkg/utils"
	http_grab_task "github.com/WangYihang/http-grab/tools/digital-ocean-docker-executor/pkg/model/task"
	"github.com/WangYihang/http-grab/tools/digital-ocean-docker-executor/pkg/option"
	"github.com/charmbracelet/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
	fd, err := os.OpenFile(option.Opt.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("failed to open log file", "error", err, "path", option.Opt.LogFilePath)
		os.Exit(1)
	}
	log.SetOutput(gojob_utils.NewTeeWriterCloser(os.Stdout, fd))
}

func main() {
	log.Info("starting", "options", option.Opt)
	s := scheduler.New(option.Opt.Name).
		WithProvider(provider.Use("digitalocean", option.Opt.DigitalOceanToken)).
		WithCreateServerOptions(
			api.NewCreateServerOptions().
				WithName(option.Opt.DropletName).
				WithTag(option.Opt.Name).
				WithRegion(option.Opt.DropletRegion).
				WithSize(option.Opt.DropletSize).
				WithImage(option.Opt.DropletImage).
				WithPrivateKeyPath(option.Opt.DropletPrivateKeyPath).
				WithPublicKeyPath(option.Opt.DropletPublicKeyPath).
				WithPublicKeyName(option.Opt.Name),
		).
		WithMaxConcurrency(option.Opt.NumDroplets).
		WithDestroyAfterFinished(true)
	for t := range http_grab_task.Generate(option.Opt.Name, option.Opt.InputFilePath, 80) {
		s.Submit(t)
	}
	s.Wait()
}
