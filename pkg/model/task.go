package model

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Task struct {
	IP                 string            `json:"ip"`
	Port               uint16            `json:"port"`
	Scheme             string            `json:"scheme"`
	Method             string            `json:"method"`
	Path               string            `json:"path"`
	Host               string            `json:"host"`
	Body               string            `json:"body"`
	Queries            map[string]string `json:"queries"`
	Headers            map[string]string `json:"headers"`
	SNI                string            `json:"sni"`
	HTTP               HTTP              `json:"http"`
	TLS                *TLS              `json:"tls,omitempty"`
	ConnectTimeout     int               `json:"connect_timeout"`
	Timeout            int               `json:"timeout"`
	InsecureSkipVerify bool              `json:"insecure_skip_verify"`
	Error              string            `json:"error"`
}

func NewTask(line string) *Task {
	ip := strings.TrimSpace(line)
	return &Task{
		IP:                 ip,
		Port:               80,
		Scheme:             "http",
		Method:             "GET",
		Path:               "/",
		SNI:                ip,
		Host:               ip,
		Queries:            make(map[string]string),
		Headers:            make(map[string]string),
		Body:               "",
		HTTP:               HTTP{},
		TLS:                nil,
		ConnectTimeout:     4,
		Timeout:            8,
		InsecureSkipVerify: false,
		Error:              "",
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
	if host == "" {
		host = t.IP
	}
	t.Host = host
	return t
}

func (t *Task) WithSNI(sni string) *Task {
	t.SNI = sni
	return t
}

func (t *Task) WithConnectTimeout(connectTimeout int) *Task {
	t.ConnectTimeout = connectTimeout
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

func (t *Task) WithHeader(key, value string) *Task {
	if t.Headers == nil {
		t.Headers = make(map[string]string)
	}
	t.Headers[key] = value
	return t
}

func (t *Task) WithQuery(key, value string) *Task {
	if t.Queries == nil {
		t.Queries = make(map[string]string)
	}
	t.Queries[key] = value
	return t
}

func (t *Task) WithBody(body string) *Task {
	t.Body = body
	return t
}

func (t *Task) Do() error {
	// Create HTTP Client
	tlsClientConfig := &tls.Config{
		ServerName:         t.SNI,
		InsecureSkipVerify: t.InsecureSkipVerify,
	}
	transport := &http.Transport{
		DisableCompression: true,
		TLSClientConfig:    tlsClientConfig,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: time.Duration(t.ConnectTimeout) * time.Second,
			}
			return dialer.DialContext(ctx, network, fmt.Sprintf("%s:%d", t.IP, t.Port))
		},
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: time.Duration(t.ConnectTimeout) * time.Second,
			}
			conn, err := tls.DialWithDialer(dialer, network, addr, tlsClientConfig)
			if err != nil {
				return nil, err
			}
			connectionState := conn.ConnectionState()
			t.TLS = NewTLS(&connectionState)
			return conn, nil
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
	u := t.URL()
	req, err := http.NewRequest(t.Method, u.String(), nil)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}
	req.Close = true
	req.Host = u.Host
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")
	req.Body = io.NopCloser(strings.NewReader(t.Body))

	// Create HTTP Request
	httpRequest, err := NewHTTPRequest(req)
	if err != nil {
		slog.Debug("error occured while creating http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}
	t.HTTP.Request = httpRequest

	// Do HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("error occured while doing http request", slog.String("error", err.Error()))
		t.Error = err.Error()
		return err
	}

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

func (t *Task) URL() *url.URL {
	u := &url.URL{
		Scheme: t.Scheme,
		Host:   t.Host,
		Path:   t.Path,
	}
	if t.Port != 0 && !isDefaultPort(t.Scheme, t.Port) {
		u.Host = fmt.Sprintf("%s:%d", t.Host, t.Port)
	}
	query := u.Query()
	for k, v := range t.Queries {
		query.Add(k, v)
	}
	u.RawQuery = query.Encode()
	return u
}

func isDefaultPort(scheme string, port uint16) bool {
	return (scheme == "http" && port == 80) || (scheme == "https" && port == 443)
}
