package model

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type Task struct {
	IP                 string `json:"ip"`
	Port               uint16 `json:"port"`
	Scheme             string `json:"scheme"`
	Method             string `json:"method"`
	Path               string `json:"path"`
	Host               string `json:"host"`
	SNI                string `json:"sni"`
	HTTP               HTTP   `json:"http"`
	Timeout            int    `json:"timeout"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	Error              string `json:"error"`
}

func NewTask(line string) *Task {
	ip := strings.TrimSpace(line)
	return &Task{
		IP:                 ip,
		Port:               80,
		Scheme:             "http",
		Method:             "GET",
		Path:               "/",
		Host:               ip,
		SNI:                ip,
		Timeout:            8,
		Error:              "",
		InsecureSkipVerify: false,
	}
}

func (t *Task) WithPort(port uint16) *Task {
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

func (t *Task) WithSNI(sni string) *Task {
	t.SNI = sni
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

func (t *Task) WithInsecureSkipVerify(insecureSkipVerify bool) *Task {
	t.InsecureSkipVerify = insecureSkipVerify
	return t
}

func (t *Task) Do() error {
	// Create HTTP Client
	dialer := &net.Dialer{}
	transport := &http.Transport{
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			ServerName:         t.SNI,
			InsecureSkipVerify: t.InsecureSkipVerify,
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Override the address with the specified IP
			return dialer.DialContext(ctx, network, fmt.Sprintf("%s:%d", t.IP, t.Port))
		},
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
