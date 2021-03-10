package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/damonchen/oss-server/internal/config"
)

type tencent struct{}

func (f tencent) Create(cfg *config.Configuration) ProxyProvider {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		cfg.Tencent.Bucket, cfg.Tencent.Region))

	log.Debugf("create tencent url %s, cfg: %s", u, cfg)

	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.Tencent.ApiID,
			SecretKey: cfg.Tencent.ApiKey,
		},
	})

	proxy := &TencentProxy{
		client: client,
	}
	return proxy
}

type TencentProxy struct {
	client *cos.Client
}

func (proxy *TencentProxy) Handle(w http.ResponseWriter, req *http.Request) {
	paths := req.URL.Query()["path"]
	if len(paths) == 0 {
		log.Fatal("tencent proxy path is empty %s", w)
		w.WriteHeader(500)
		return
	}

	path := strings.TrimSpace(paths[0])
	if len(path) == 0 {
		log.Fatal("tencent proxy path is empty string")
		w.WriteHeader(500)
		return
	}

	c := proxy.client
	resp, err := c.Object.Get(context.Background(), path, nil)
	if err != nil {
		log.Fatalf("object get from tencent oss error %+v", err)
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
