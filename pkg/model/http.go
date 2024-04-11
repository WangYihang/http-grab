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

type HTTPRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Host   string `json:"host"`

	RemoteAddr string `json:"remote_addr"`
	RequestURI string `json:"request_uri"`

	Proto      string `json:"proto"`
	ProtoMajor int    `json:"proto_major"`
	ProtoMinor int    `json:"proto_minor"`

	Headers http.Header `json:"headers"`

	ContentLength    int64    `json:"content_length"`
	TransferEncoding []string `json:"transfer_encoding,omitempty"`
	Close            bool     `json:"close"`

	Form          url.Values      `json:"form,omitempty"`
	PostForm      url.Values      `json:"post_form,omitempty"`
	MultipartForm *multipart.Form `json:"multipart_form,omitempty"`

	RawBody    []byte `json:"raw_body"`
	BodySha256 string `json:"body_sha256"`
	Body       string `json:"body"`

	Trailer http.Header `json:"trailer,omitempty"`
}

func NewHTTPRequest(req *http.Request) (*HTTPRequest, error) {
	rawBody := []byte{}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			slog.Warn("error occured while reading request body", slog.String("error", err.Error()))
			return nil, err
		}
		rawBody = body
	}
	req.Body = io.NopCloser(bytes.NewReader(rawBody))
	httpRequest := &HTTPRequest{
		Method:           req.Method,
		URL:              req.URL.String(),
		Host:             req.Host,
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Headers:          req.Header,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Close:            req.Close,
		Form:             req.Form,
		PostForm:         req.PostForm,
		MultipartForm:    req.MultipartForm,
		Trailer:          req.Trailer,
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

	Headers http.Header `json:"headers"`

	RawBody    []byte `json:"raw_body"`
	BodySha256 string `json:"body_sha256"`
	Body       string `json:"body"`

	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding,omitempty"`
	Close            bool        `json:"close"`
	Uncompressed     bool        `json:"uncompressed"`
	Trailer          http.Header `json:"trailer"`
}

func NewHTTPResponse(resp *http.Response) (*HTTPResponse, error) {
	httpResponse := &HTTPResponse{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Headers:          resp.Header,
		RawBody:          []byte{},
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          resp.Trailer,
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("error occured while reading response body", slog.String("error", err.Error()))
		return httpResponse, nil
	}
	httpResponse.RawBody = body
	httpResponse.Body = string(body)
	httpResponse.BodySha256 = util.Sha256(body)
	return httpResponse, nil
}
