package web

import (
	"fmt"
	"github.com/damonchen/oss-server/internal/web/provider"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/damonchen/oss-server/internal/config"
)

//// Example format string. Everything except the message has a custom color
//// which is dependent on the log level. Many fields have a custom output
//// formatting too, eg. the time returns the hour down to the milli second.
//var format = logging.MustStringFormatter(
//	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
//)

var log = logging.MustGetLogger("web")

type Server struct {
	Cfg *config.Configuration
}

func (svr *Server) Run() error {
	// create provider server client
	for _, pvd := range svr.Cfg.Providers {
		if !provider.IsSupportProvider(pvd) {
			log.Fatalf("not support provider %s provide", pvd)
			return errors.Errorf("not support provider %s", pvd)
		}

		factory := provider.GetFactory(pvd)
		if factory == nil {
			log.Fatalf("not support provider %s provide", pvd)
			return errors.Errorf("not support provider %s", pvd)
		}
		proxyProvider := factory.Create(svr.Cfg)
		provider.RegisterProxyProvider(pvd, proxyProvider)
	}

	r := NewRouter(svr)

	log.Infof("start listen port: %s", svr.Cfg.Port)

	port := fmt.Sprintf(":%s", svr.Cfg.Port)
	return http.ListenAndServe(port, r)
}

func (svr *Server) Handle(w http.ResponseWriter, req *http.Request) {
	pvd := chi.URLParam(req, "provider")
	if pvd == "" {
		log.Fatalf("provider is empty")
		w.WriteHeader(500)
	}

	if !provider.IsSupportProvider(pvd) {
		log.Fatalf("not support provider %s provide", pvd)
		w.WriteHeader(500)
	}

	proxyProvider := provider.GetProxyProvider(pvd)
	proxyProvider.Handle(w, req)
}

func (svr *Server) Auth(next http.Handler) http.Handler {
	auth := svr.Cfg.Auth
	var isNoneAuth = auth.Type == "none" || auth.Type == ""
	var isBasicAuth = auth.Type == "basic"
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 判定下是否已经认证通过了
		if isNoneAuth {
			next.ServeHTTP(w, req)
			return
		}

		if isBasicAuth {
			name, password, ok := req.BasicAuth()
			if !ok {
				log.Fatalf("not supply correct basic auth")
				w.WriteHeader(500)
				return
			}

			if name == auth.BasicAuth.Name && password == auth.BasicAuth.Password {
				next.ServeHTTP(w, req)
			} else {
				log.Debugf("basic name %s and password %s not correct", name, password)
				w.WriteHeader(401)
			}
		}

		err := proxyAuth(auth.ProxyAuth, w, req)
		if err != nil {
			log.Errorf("proxy auth failed: %s", err)
			// 401错误
		}

		next.ServeHTTP(w, req)
		return
	})
}
