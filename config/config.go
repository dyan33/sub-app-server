package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var Cfg *Config

type Config struct {
	WebPort      int    `yaml:"web_port"`
	ProxyPort    int    `yaml:"proxy_port"`
	ProxyNum     int    `yaml:"proxy_num"`
	Python       string `yaml:"python"`
	OrangeScript string `yaml:"orange_script"`
}

func init() {

	log.SetOutput(os.Stdout)

	Cfg = &Config{}

	if data, err := ioutil.ReadFile("config.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), Cfg)
	}
	log.Println("加载配置", Cfg)
}
