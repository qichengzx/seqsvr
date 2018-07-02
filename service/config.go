package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	PORT  string `yaml:"port"`
	STEP  int64  `yaml:"step"`
	MySQL MySQL  `yaml:"mysql"`
}

type MySQL struct {
	User     string `yaml:"user"`
	PassWord string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	MaxIdle  int    `yaml:"maxidle"`
	MaxOpen  int    `yaml:"maxopen"`
}

func NewConfig() *Config {
	b, _ := ioutil.ReadFile("config.yml")
	conf := new(Config)
	yaml.Unmarshal(b, conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return conf
}
