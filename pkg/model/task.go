package model

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Task struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Host    string `json:"host"`
	HTTP    HTTP   `json:"http"`
	Timeout int    `json:"timeout"`
}

func NewTask(line string, port int, path string, host string, timeout int) *Task {
	ip := strings.TrimSpace(line)
	if host == "" {
		host = ip
	}
	return &Task{
		IP:      ip,
		Port:    port,
		Path:    path,
		Host:    host,
		Timeout: timeout,
	}
}

func (t *Task) Do() error {
	// Create HTTP Client
	transport := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout:   time.Duration(t.Timeout) * time.Second,
		Transport: transport,
	}

	// Create HTTP Request
	u := fmt.Sprintf("http://%s:%d/%s", t.IP, t.Port, t.Path)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		return err
	}
	req.Host = t.Host
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")

	// Do HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("error occured while doing http request", slog.String("error", err.Error()))
		return err
	}

	// Create HTTP Request
	httpRequest, err := NewHTTPRequest(req)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		return err
	}
	t.HTTP.Request = httpRequest

	// Create HTTP Response
	httpResponse, err := NewHTTPResponse(resp)
	if err != nil {
		slog.Debug("error occured while creating http response", slog.String("error", err.Error()))
		return err
	}
	t.HTTP.Response = httpResponse
	return nil
}
