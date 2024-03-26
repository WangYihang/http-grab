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
	Scheme  string `json:"scheme"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Host    string `json:"host"`
	HTTP    HTTP   `json:"http"`
	Timeout int    `json:"timeout"`
	Error   string `json:"error"`
}

func NewTask(line string) *Task {
	ip := strings.TrimSpace(line)
	return &Task{
		IP:      ip,
		Port:    80,
		Scheme:  "http",
		Method:  "GET",
		Path:    "/",
		Host:    ip,
		Timeout: 8,
		Error:   "",
	}
}

func (t *Task) WithPort(port int) *Task {
	t.Port = port
	return t
}

func (t *Task) WithPath(path string) *Task {
	t.Path = path
	return t
}

func (t *Task) WithHost(host string) *Task {
	t.Host = host
	return t
}

func (t *Task) WithTimeout(timeout int) *Task {
	t.Timeout = timeout
	return t
}

func (t *Task) WithScheme(scheme string) *Task {
	t.Scheme = scheme
	return t
}

func (t *Task) WithMethod(method string) *Task {
	t.Method = method
	return t
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
	if t.Path == "" {
		t.Path = "/"
	}
	u := fmt.Sprintf("%s://%s:%d%s", t.Scheme, t.IP, t.Port, t.Path)
	req, err := http.NewRequest(t.Method, u, nil)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}
	req.Host = t.Host
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")

	// Do HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("error occured while doing http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}

	// Create HTTP Request
	httpRequest, err := NewHTTPRequest(req)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}
	t.HTTP.Request = httpRequest

	// Create HTTP Response
	httpResponse, err := NewHTTPResponse(resp)
	if err != nil {
		slog.Debug("error occured while creating http response", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}
	t.HTTP.Response = httpResponse
	return nil
}
