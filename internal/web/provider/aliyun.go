package provider

import (
	"fmt"
	"github.com/damonchen/oss-server/internal/web/utils"
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
	path, err := utils.GetRequestPath(req)
	if err != nil {
		log.Errorf("request path error %s", err)
		w.WriteHeader(500)
		return
	}
	path = strings.TrimPrefix(path, "/")

	bucket, err := proxy.client.Bucket(proxy.bucket)
	if err != nil {
		log.Errorf("get bucket %s error %s", proxy.bucket, err)
		w.WriteHeader(500)
	}

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
