package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func Load(fileName string, cfg *Configuration) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, cfg)
}
