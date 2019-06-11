package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var C *Config

type Cache struct {
	Dir    string   `json:"dir"`
	Expire string   `json:"expire"`
	Types  []string `json:"types"`
	Urls   []string `json:"urls"`
	Ignore []string `json:"ignore"`
}

type Config struct {
	Server []int    `yaml:"server"`
	Proxy  []int    `yaml:"proxy"`
	Igonre []string `yaml:"igonre"`

	Cache Cache `yaml:"cache"`
}

func init() {

	log.SetOutput(os.Stdout)

	C = &Config{}

	if data, err := ioutil.ReadFile("config.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), C)
	}

	log.Println("加载配置", C)
}
