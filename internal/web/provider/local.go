package provider

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"

	"github.com/damonchen/oss-server/internal/config"
	"github.com/damonchen/oss-server/internal/web/handle"
	"github.com/damonchen/oss-server/internal/web/utils"
)

type local struct{}

func (l local) Create(cfg *config.Configuration) []ProxyProvider {
	var proxyProviders []ProxyProvider
	for _, proxy := range cfg.Local {

		proxy := &LocalProxy{
			urlPrefix: proxy.UrlPrefix,
			name:      proxy.Name,
			path:      proxy.Path,
		}

		proxyProviders = append(proxyProviders, proxy)
	}

	return proxyProviders
}

type LocalProxy struct {
	urlPrefix string
	name      string
	path      string
}

func (proxy *LocalProxy) Handle(w http.ResponseWriter, req *http.Request) {
	path, err := utils.GetRequestPath(req)
	if err != nil {
		log.Errorf("request path error %s", err)
		w.WriteHeader(500)
		return
	}
	path = strings.TrimPrefix(path, "/")
	log.Debugf("proxy path file handle %s", path)

	dir := path[:2]
	dir = filepath.Join(proxy.path, dir)

	filename := filepath.Join(dir, path[2:])

	log.Debugf("will response file %s", filename)

	fp, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		log.Errorf("read file path %s error %s", filename, err)
		w.WriteHeader(500)
		return
	}

	data, err := io.ReadAll(fp)
	if err != nil {
		log.Errorf("read file error %s", err)
		w.WriteHeader(500)
		return
	}

	buff := data[:100]
	kind, _ := filetype.Match(buff)
	if kind == filetype.Unknown {
		w.Write(data)
	} else {
		w.Header().Add("Content-Type", kind.MIME.Value)
		w.Write(data)
	}

}

func (proxy *LocalProxy) Upload(w http.ResponseWriter, req *http.Request) {
	// 本地路径的上传方式
	file, _, err := req.FormFile("file")
	if err != nil {
		log.Errorf("request get file error %s", err)
		handle.HandleError(w, err)
		return
	}

	// ext := filepath.Ext(header.Filename)
	randomString := utils.GetRandomString(64)
	filename := utils.GetFilename(randomString)

	log.Debugf("get filename %s", filename)

	dir := filename[:2]
	dir = filepath.Join(proxy.path, dir)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				log.Errorf("mkdir %s error", dir, err)
				handle.HandleError(w, err)
				return
			}
		}
	}

	newFilename := filename[2:]
	newFilename = filepath.Join(dir, newFilename)

	log.Debugf("will save file %s", newFilename)

	fp, err := os.OpenFile(newFilename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Errorf("open file %s error", newFilename, err)
		handle.HandleError(w, err)
		return
	}

	var buff bytes.Buffer
	reader := io.TeeReader(file, &buff)
	_, err = io.Copy(fp, reader)
	if err != nil {
		log.Errorf("save file %s error", newFilename, err)
		handle.HandleError(w, err)
		return
	}

	// TODO: 图片合规性检查

	data := map[string]interface{}{}
	data["url"] = filepath.Join(proxy.urlPrefix, filename)

	handle.HandleSuccess(w, data)
}

func (proxy *LocalProxy) Name() string {
	return proxy.name
}

func init() {
	RegisterFactory("local", &local{})
}
