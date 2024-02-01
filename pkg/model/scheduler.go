package model

import (
	"bufio"
	"log/slog"
	"os"
	"sync"
)

func FanIn[T any](cs []chan T) chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c chan T) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FanOut[T any](tasks chan T, n int) []chan T {
	outChans := make([]chan T, n)
	for i := range outChans {
		outChans[i] = make(chan T)
	}
	go func() {
		defer func() {
			for _, c := range outChans {
				close(c)
			}
		}()
		var wg sync.WaitGroup
		wg.Add(n)
		for _, c := range outChans {
			go func(c chan T) {
				defer wg.Done()
				for task := range tasks {
					c <- task
				}
			}(c)
		}
		wg.Wait()
	}()
	return outChans
}

type ITask interface {
	Do()
	JSON() string
}

type TaskFactory[T ITask] func(index int, ip string, port int, path string, host string, timeout int) T

func LoadTasks[T ITask](factory TaskFactory[T], inputFilePath string, port int, path string, host string, timeout int) chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		index := 0
		for line := range ReadFile(inputFilePath) {
			task := factory(index, line, port, path, host, timeout)
			out <- task
			index++
		}
	}()
	return out
}

func StoreTasks[T ITask](tasks chan T, path string) {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("error occured while opening file", slog.String("error", err.Error()))
		return
	}
	defer fd.Close()
	for task := range tasks {
		fd.WriteString(task.JSON() + "\n")
	}
}

func Worker[T ITask](tasks chan T) chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for task := range tasks {
			task.Do()
			out <- task
		}
	}()
	return out
}

func ReadFile(path string) chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		fd, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			slog.Error("error occured while opening file", slog.String("error", err.Error()))
			return
		}
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			slog.Error("error occured while scanning file", slog.String("error", err.Error()))
			return
		}
		fd.Close()
	}()
	return out
}
