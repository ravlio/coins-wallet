package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// load config from yaml file into the struct
func Load(path string, dst interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read config file error")
	}

	return LoadFromBytes(b, dst)
}

// load config from byte slice
func LoadFromBytes(data []byte, dst interface{}) error {
	err := yaml.Unmarshal([]byte(data), dst)

	if err != nil {
		return errors.Wrap(err, "yaml unmarshal error")
	}

	return nil
}
