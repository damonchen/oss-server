package config

import (
	"github.com/ghodss/yaml"
	"github.com/op/go-logging"
	"io/ioutil"
)

var (
	log = logging.MustGetLogger("config")
)

func Load(fileName string) (*Configuration, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("load config file %s error %s", fileName, err)
		return nil, err
	}

	cfg := Configuration{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Errorf("load yaml config error %s", err)
		return nil, err
	}

	// TODO: default config
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	return &cfg, nil
}
