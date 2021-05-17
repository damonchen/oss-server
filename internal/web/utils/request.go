package utils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
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

func GetRequestType(req *http.Request) (string, error) {
	typs := req.URL.Query()["type"]
	if len(typs) == 0 {
		err := errors.New("proxy path is empty")
		return "", err
	}

	typ := strings.TrimSpace(typs[0])
	if len(typ) == 0 {
		err := errors.New("proxy type is empty string")
		return "", err
	}

	// 对typ做url decode的处理
	typ, err := url.QueryUnescape(typ)
	if err != nil {
		return "", err
	}

	return typ, nil
}
