package main

import (
	"os"

	"github.com/WangYihang/http-grab/pkg/model"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFilePath  string `short:"i" long:"input" description:"input file path" required:"true"`
	OutputFilePath string `short:"o" long:"output" description:"output file path" required:"true"`
	NumWorkers     int    `short:"n" long:"num-workers" description:"number of workers" default:"32"`
	Timeout        int    `short:"t" long:"timeout" description:"timeout" default:"8"`
	Port           int    `short:"p" long:"port" description:"port" default:"80"`
	Path           string `short:"P" long:"path" description:"path" default:"index.html"`
	Host           string `short:"H" long:"host" description:"host" default:""`
}

var opts Options

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	numWorkers := opts.NumWorkers
	resultChans := make([]chan *model.Task, 0, numWorkers)
	for _, taskChan := range model.FanOut(model.LoadTasks(
		model.NewTask,
		opts.InputFilePath,
		opts.Port,
		opts.Path,
		opts.Host,
		opts.Timeout,
	), numWorkers) {
		resultChans = append(resultChans, model.Worker(taskChan))
	}
	model.StoreTasks(model.FanIn(resultChans), opts.OutputFilePath)
}
