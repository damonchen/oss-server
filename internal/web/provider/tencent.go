package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/damonchen/oss-server/internal/web/utils"
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
			name:             tencent.Name,
			client:           client,
			defaultImagePath: tencent.DefaultImagePath,
		}

		proxyProviders = append(proxyProviders, proxy)
	}

	return proxyProviders
}

type TencentProxy struct {
	name             string
	client           *cos.Client
	defaultImagePath string
}

func (proxy *TencentProxy) Name() string {
	return proxy.name
}

func (proxy *TencentProxy) Handle(w http.ResponseWriter, req *http.Request) {
	path, err := utils.GetRequestPath(req)
	if err != nil {
		log.Errorf("request path error %s", err)
		w.WriteHeader(500)
		return
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	log.Debugf("will get tencent path: %+v", path)

	c := proxy.client
	resp, err := c.Object.Get(context.Background(), path, nil)
	if err != nil {
		// 当没有找到的时候，获取参数类型，然后返回默认的大小
		typ, err := utils.GetRequestType(req)
		if err != nil {
			log.Errorf("object get from tencent oss error %+v", err)
			w.WriteHeader(404)
			return
		}
		typ = strings.ToLower(typ)
		log.Debugf("tencent type %s", typ)
		// 返回2倍图或者3倍图
		var filename string
		if typ == "icon" {
			filename = filepath.Join(proxy.defaultImagePath, "user_icon_nor@2x.png")
		} else if typ == "avatar" {
			filename = filepath.Join(proxy.defaultImagePath, "user_head_nor@2x.png")
		} else if typ == "head" {
			filename = filepath.Join(proxy.defaultImagePath, "home_nor@2x.png")
		} else {
			log.Errorf("object get from tencent oss error %+v", err)
			w.WriteHeader(404)
			return
		}

		fp, err := os.Open(filename)
		if err != nil {
			log.Errorf("open file %s error %s", filename, err)
			w.WriteHeader(500)
			return
		}
		defer fp.Close()
		_, _ = io.Copy(w, fp)
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
}

func (proxy *TencentProxy) Upload(w http.ResponseWriter, req *http.Request) {
	// TODO:
}

func init() {
	RegisterFactory("tencent", &tencent{})
}
