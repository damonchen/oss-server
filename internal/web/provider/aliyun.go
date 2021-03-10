package provider

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/damonchen/oss-server/internal/config"
)

type aliyun struct{}

func (f aliyun) Create(cfg *config.Configuration) []ProxyProvider {

	var proxyProviders []ProxyProvider
	for _, proxy := range cfg.Aliyun {
		endpoint := fmt.Sprintf("%s.aliyuncs.com", proxy.Region)
		client, err := oss.New(endpoint, proxy.ApiID, proxy.ApiKey)
		if err != nil {
			return nil
		}

		proxy := &AliyunProxy{
			client: client,
			bucket: proxy.Bucket,
			name:   proxy.Name,
		}

		proxyProviders = append(proxyProviders, proxy)
	}

	return proxyProviders
}

type AliyunProxy struct {
	client *oss.Client
	bucket string
	name   string
}

func (proxy *AliyunProxy) Handle(w http.ResponseWriter, req *http.Request) {
	bucket, err := proxy.client.Bucket(proxy.bucket)
	if err != nil {
		log.Errorf("get bucket %s error %s", proxy.bucket, err)
		w.WriteHeader(500)
	}

	paths := req.URL.Query()["path"]
	if len(paths) == 0 {
		log.Error("aliyun proxy path is empty")
		w.WriteHeader(500)
		return
	}

	path := strings.TrimSpace(paths[0])
	if len(path) == 0 {
		log.Error("aliyun proxy path is empty string")
		w.WriteHeader(500)
		return
	}

	path = strings.TrimPrefix(path, "/")
	body, err := bucket.GetObject(path)
	if err != nil {
		log.Errorf("download bucket %s file %s error %s", proxy.bucket, path, err)
		w.WriteHeader(500)
		return
	}
	defer body.Close()

	_, _ = io.Copy(w, body)
	return
}

func (proxy *AliyunProxy) Name() string {
	return proxy.name
}

func init() {
	RegisterFactory("aliyun", &aliyun{})
}
