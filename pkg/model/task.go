package model

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Task struct {
	Index      int    `json:"index"`
	StartedAt  int64  `json:"started_at"`
	FinishedAt int64  `json:"finished_at"`
	NumTries   int    `json:"num_tries"`
	Timeout    int    `json:"timeout"`
	Error      string `json:"error"`

	IP   string `json:"ip"`
	Port int    `json:"port"`
	Path string `json:"path"`
	Host string `json:"host"`
	HTTP HTTP   `json:"http"`
}

func NewTask(index int, line string, port int, path string, host string, timeout int) *Task {
	ip := strings.TrimSpace(line)
	if host == "" {
		host = ip
	}
	return &Task{
		Index:   index,
		IP:      ip,
		Port:    80,
		Path:    path,
		Timeout: timeout,
		Host:    host,
	}
}

func (t *Task) Do() {
	t.StartedAt = time.Now().UnixMilli()
	defer func() {
		t.FinishedAt = time.Now().UnixMilli()
	}()
	for i := 0; i < 1; i++ {
		// Increase NumTries
		t.NumTries++

		// Create HTTP Client
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 8 * time.Second,
		}

		// Create HTTP Request
		u := fmt.Sprintf("http://%s:%d/%s", t.IP, t.Port, t.Path)
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			t.Error = err.Error()
			slog.Error("error occured while creating http request", slog.String("error", err.Error()))
			continue
		}
		req.Host = t.Host
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")

		// Do HTTP Request
		resp, err := client.Do(req)
		if err != nil {
			t.Error = err.Error()
			slog.Error("error occured while doing http request", slog.String("error", err.Error()))
			continue
		}

		// Create HTTP Request
		httpRequest, err := NewHTTPRequest(req)
		if err != nil {
			t.Error = err.Error()
			slog.Error("error occured while creating http request", slog.String("error", err.Error()))
			continue
		}
		t.HTTP.Request = httpRequest

		// Create HTTP Response
		httpResponse, err := NewHTTPResponse(resp)
		if err != nil {
			t.Error = err.Error()
			slog.Error("error occured while creating http response", slog.String("error", err.Error()))
			continue
		}
		t.HTTP.Response = httpResponse
		break
	}
}

func (t *Task) JSON() string {
	data, err := json.Marshal(t)
	if err != nil {
		slog.Error("error occured while marshalling task", slog.String("error", err.Error()))
		return ""
	}
	return string(data)
}
