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

func (f tencent) Create(cfg *config.Configuration) []ProxyProvider {
	var proxyProviders []ProxyProvider
	for _, tencent := range cfg.Tencent {
		u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
			tencent.Bucket, tencent.Region))

		log.Debugf("create tencent %s url %s, ", tencent.Name, u)

		b := &cos.BaseURL{BucketURL: u}

		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  tencent.ApiID,
				SecretKey: tencent.ApiKey,
			},
		})

		proxy := &TencentProxy{
			name:   tencent.Name,
			client: client,
		}

		proxyProviders = append(proxyProviders, proxy)
	}

	return proxyProviders
}

type TencentProxy struct {
	name   string
	client *cos.Client
}

func (proxy *TencentProxy) Name() string {
	return proxy.name
}

func (proxy *TencentProxy) Handle(w http.ResponseWriter, req *http.Request) {
	paths := req.URL.Query()["path"]
	if len(paths) == 0 {
		log.Error("tencent proxy path is empty")
		w.WriteHeader(500)
		return
	}

	path := strings.TrimSpace(paths[0])
	if len(path) == 0 {
		log.Error("tencent proxy path is empty string")
		w.WriteHeader(500)
		return
	}

	c := proxy.client
	resp, err := c.Object.Get(context.Background(), path, nil)
	if err != nil {
		log.Errorf("object get from tencent oss error %+v", err)
		w.WriteHeader(500)
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
}

func init() {
	RegisterFactory("tencent", &tencent{})
}
