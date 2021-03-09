package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/damonchen/oss-server/internal/config"
)

type authResp struct {
	// status must success or fail
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func proxyAuth(proxyAuth config.ProxyAuth, w http.ResponseWriter, req *http.Request) error {
	client := resty.New()

	headers := map[string]string{}
	for _, keys := range req.Header {
		key := keys[0]
		var v = req.Header.Get(key)
		headers[key] = v
	}

	resp, err := client.R().SetHeaders(headers).Get(proxyAuth.AuthPath)
	if err != nil {
		return err
	}

	b := resp.Body()
	ar := authResp{}
	err = json.Unmarshal(b, &ar)
	if err != nil {
		return err
	}

	if !ar.Status {
		return errors.Errorf("auth proxy path failed: %s", ar.Message)
	}
	return nil
}
