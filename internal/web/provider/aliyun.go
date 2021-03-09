package provider

import (
	"net/http"

	"github.com/damonchen/oss-server/internal/config"
)

type aliyun struct{}

func (f aliyun) Create(cfg *config.Configuration) ProxyProvider {
	proxy := AliyunProxy{

	}
	return proxy
}

type AliyunProxy struct {
}

func (proxy AliyunProxy) Handle(resp http.ResponseWriter, req *http.Request) {
	// TODO
}
