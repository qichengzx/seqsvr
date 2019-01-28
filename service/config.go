package service

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	PORT  string `yaml:"port"`
	STEP  uint64 `yaml:"step"`
	MySQL MySQL  `yaml:"mysql"`
}

func NewConfig() *Config {
	b, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	conf := new(Config)
	yaml.Unmarshal(b, conf)

	return conf
}
