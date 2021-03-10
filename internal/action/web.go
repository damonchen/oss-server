package action

import (
	"github.com/damonchen/oss-server/internal/config"
	web2 "github.com/damonchen/oss-server/internal/web"
)

type Web struct {
	Config  *config.Configuration
	CfgFile string
}

func (web *Web) Run() error {
	cfg, err := config.Load(web.CfgFile)
	if err != nil {
		return err
	}

	web.Config = cfg

	return web.serverRun()
}

func (web *Web) serverRun() error {
	srv := web2.Server{Cfg: web.Config}
	return srv.Run()
}
