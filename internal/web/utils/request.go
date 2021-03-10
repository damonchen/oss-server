package utils

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

func GetRequestPath(req *http.Request) (string, error) {
	paths := req.URL.Query()["path"]
	if len(paths) == 0 {
		err := errors.New("proxy path is empty")
		return "", err
	}

	path := strings.TrimSpace(paths[0])
	if len(path) == 0 {
		err := errors.New("proxy path is empty string")
		return "", err
	}

	// 对path做url decode的处理
	path, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
