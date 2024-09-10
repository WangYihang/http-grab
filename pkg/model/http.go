package model

import (
	"bytes"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/WangYihang/http-grab/pkg/util"
)

type HTTP struct {
	Request  *HTTPRequest  `json:"request,omitempty"`
	Response *HTTPResponse `json:"response,omitempty"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HTTPRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Host   string `json:"host"`

	RemoteAddr string `json:"remote_addr"`
	RequestURI string `json:"request_uri"`

	Path string `json:"path"`

	Proto      string `json:"proto"`
	ProtoMajor int    `json:"proto_major"`
	ProtoMinor int    `json:"proto_minor"`

	Headers []Header `json:"headers"`

	ContentLength    int64    `json:"content_length"`
	TransferEncoding []string `json:"transfer_encoding,omitempty"`
	Close            bool     `json:"close"`

	Form          url.Values      `json:"form,omitempty"`
	PostForm      url.Values      `json:"post_form,omitempty"`
	MultipartForm *multipart.Form `json:"multipart_form,omitempty"`

	RawBody    []byte `json:"raw_body"`
	BodySha256 string `json:"body_sha256"`
	Body       string `json:"body"`

	Trailer []Header `json:"trailer,omitempty"`
}

func NewHTTPRequest(req *http.Request) (*HTTPRequest, error) {
	rawBody := []byte{}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			slog.Warn("error occurred while reading request body", slog.String("error", err.Error()))
			return nil, err
		}
		rawBody = body
	}
	req.Body = io.NopCloser(bytes.NewReader(rawBody))
	headers := []Header{}
	for name, values := range req.Header {
		for _, value := range values {
			headers = append(headers, Header{Name: name, Value: value})
		}
	}
	trailer := []Header{}
	for name, values := range req.Trailer {
		for _, value := range values {
			trailer = append(trailer, Header{Name: name, Value: value})
		}
	}
	httpRequest := &HTTPRequest{
		Method:           req.Method,
		URL:              req.URL.String(),
		Path:             req.URL.Path,
		Host:             req.Host,
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Headers:          headers,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Close:            req.Close,
		Form:             req.Form,
		PostForm:         req.PostForm,
		MultipartForm:    req.MultipartForm,
		Trailer:          trailer,
		RawBody:          rawBody,
		BodySha256:       util.Sha256(rawBody),
		Body:             string(rawBody),
	}
	return httpRequest, nil
}

type HTTPResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`

	Proto      string `json:"proto"`
	ProtoMajor int    `json:"proto_major"`
	ProtoMinor int    `json:"proto_minor"`

	Headers []Header `json:"headers"`

	RawBody    []byte `json:"raw_body"`
	BodySha256 string `json:"body_sha256"`
	Body       string `json:"body"`

	ContentLength    int64    `json:"content_length"`
	TransferEncoding []string `json:"transfer_encoding,omitempty"`
	Close            bool     `json:"close"`
	Uncompressed     bool     `json:"uncompressed"`
	Trailer          []Header `json:"trailer"`
}

func NewHTTPResponse(resp *http.Response) (*HTTPResponse, error) {
	headers := []Header{}
	for name, values := range resp.Header {
		for _, value := range values {
			headers = append(headers, Header{Name: name, Value: value})
		}
	}
	trailer := []Header{}
	for name, values := range resp.Trailer {
		for _, value := range values {
			trailer = append(trailer, Header{Name: name, Value: value})
		}
	}
	httpResponse := &HTTPResponse{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Headers:          headers,
		RawBody:          []byte{},
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          trailer,
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("error occurred while reading response body", slog.String("error", err.Error()))
		return httpResponse, nil
	}
	httpResponse.RawBody = body
	httpResponse.Body = string(body)
	httpResponse.BodySha256 = util.Sha256(body)
	return httpResponse, nil
}
