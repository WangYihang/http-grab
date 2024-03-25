package model

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

var Runner *IPInfo
var once sync.Once

func init() {
	once.Do(func() {
		ipInfo, err := getIPInfo()
		if err != nil {
			slog.Error("error occured while getting runner ip info", slog.String("error", err.Error()))
			return
		}
		Runner = ipInfo
	})
}

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func getIPInfo() (*IPInfo, error) {
	ipInfo := &IPInfo{}
	req, err := http.Get("https://ipinfo.io/")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return nil, err
	}
	return ipInfo, nil
}
