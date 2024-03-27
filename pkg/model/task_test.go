package model_test

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/WangYihang/http-grab/pkg/model"
)

func TestCloudflareTask(t *testing.T) {
	// $ curl https://www.cloudflare.com/favicon.ico --resolve www.cloudflare.com:443:104.21.75.53 -vv
	// $ wget https://www.cloudflare.com/favicon.ico
	// $ md5sum favicon.ico
	// 112ad5f84433e5f46d607f73fb64bd60  favicon.ico
	task := model.NewTask("104.21.75.53").
		WithScheme("https").
		WithPort(443).
		WithPath("/favicon.ico").
		WithHost("www.cloudflare.com").
		WithSNI("www.cloudflare.com").
		WithInsecureSkipVerify(true)
	err := task.Do()
	if err != nil {
		slog.Error("error occured while doing task", "error", err.Error())
		return
	}
	data, err := json.Marshal(task)
	if err != nil {
		slog.Error("error occured while marshalling task", "error", err.Error())
		return
	}
	hash := md5.Sum(task.HTTP.Response.RawBody)
	md5String := hex.EncodeToString(hash[:])
	slog.Info("task done", "task", string(data))
	if md5String != "112ad5f84433e5f46d607f73fb64bd60" {
		t.Errorf("md5 hash not equal")
	}
}

func TestInvalidCertificateTask(t *testing.T) {
	task := model.NewTask("104.196.220.231").
		WithScheme("https").
		WithPort(443).
		WithPath("/favicon.ico").
		WithHost("www.cloudflare.com").
		WithSNI("www.cloudflare.com").
		WithInsecureSkipVerify(true)
	err := task.Do()
	if err != nil {
		slog.Error("error occured while doing task", "error", err.Error())
		return
	}
	data, err := json.Marshal(task)
	if err != nil {
		slog.Error("error occured while marshalling task", "error", err.Error())
		return
	}
	hash := md5.Sum(task.HTTP.Response.RawBody)
	md5String := hex.EncodeToString(hash[:])
	slog.Info("task done", "task", string(data), "body_md5", md5String)
}
