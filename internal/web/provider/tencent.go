package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/damonchen/oss-server/internal/config"
)

type tencent struct{}

func (f tencent) Create(cfg *config.Configuration) ProxyProvider {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		cfg.TencentOSS.BucketName, cfg.TencentOSS.Region))
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.TencentOSS.SecretID,
			SecretKey: cfg.TencentOSS.SecretKey,
		},
	})

	proxy := TencentProxy{
		client: client,
	}
	return proxy
}

type TencentProxy struct {
	client *cos.Client
}

func (proxy TencentProxy) Handle(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query()["path"]
	if path == nil {
		w.WriteHeader(500)
		return
	}

	c := proxy.client
	resp, err := c.Object.Get(context.Background(), path[0], nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
}

func init() {
	t := &tencent{}
	RegisterFactory("tencent", t)
}
