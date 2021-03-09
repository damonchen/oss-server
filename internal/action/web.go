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
	if web.Config == nil {
		web.Config = &config.Configuration{}
	}
	err := config.Load(web.CfgFile, web.Config)
	if err != nil {
		return err
	}

	return web.serverRun()
}

func (web *Web) serverRun() error {
	srv := web2.Server{Cfg: web.Config}
	return srv.Run()
}
