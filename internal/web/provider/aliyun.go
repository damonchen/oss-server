package provider

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/damonchen/oss-server/internal/config"
	"github.com/damonchen/oss-server/internal/web/utils"
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
			client:   client,
			endpoint: proxy.Endpoint,
			appId:    proxy.ApiID,
			appKey:   proxy.ApiKey,
			region:   proxy.Region,
			bucket:   proxy.Bucket,
			name:     proxy.Name,
		}

		proxyProviders = append(proxyProviders, proxy)
	}

	return proxyProviders
}

type AliyunProxy struct {
	client   *oss.Client
	endpoint string
	appId    string
	appKey   string
	region   string
	bucket   string
	name     string
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
		w.WriteHeader(404)
		return
	}
	defer body.Close()

	_, _ = io.Copy(w, body)

}

func (proxy *AliyunProxy) Upload(w http.ResponseWriter, req *http.Request) {
	endpoint := proxy.endpoint
	appId := proxy.appId
	appKey := proxy.appKey

	client, err := oss.New(endpoint, appId, appKey)
	if err != nil {
		// return nil, mistake.New500LogicErrorOnlyErr(err)
		return
	}
	bucket, err := client.Bucket(proxy.bucket)
	if err != nil {
		// return nil, mistake.New500LogicErrorOnlyErr(err)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		return
	}

	// TODO 依据不同的业务（req.Kind）将文件上传至不同位置
	var buf bytes.Buffer

	ext := filepath.Ext(header.Filename)
	fmt.Println(ext)
	// village := middleware.MustVillageFromContext(l.ctx)
	// objectKey := fmt.Sprintf("%s/%s%s", village.Name.String, uuid.NewV4().String(), )
	objectKey := ""
	if err = bucket.PutObject(objectKey, io.TeeReader(file, &buf)); err != nil {
		log.Errorf("Upload file %q error: %s", objectKey, err)
		// return nil, mistake.New500LogicErrorOnlyErr(err)
	}
}

func (proxy *AliyunProxy) Name() string {
	return proxy.name
}

func init() {
	RegisterFactory("aliyun", &aliyun{})
}
