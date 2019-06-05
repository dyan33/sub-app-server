package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var C *Config

type Cache struct {
	Types []string `json:"types"`
	Urls  []string `json:"urls"`
}

type Config struct {
	Server []int    `yaml:"server"`
	Proxy  []int    `yaml:"proxy"`
	Igonre []string `yaml:"igonre"`

	Cache Cache `yaml:"cache"`

	CacheDir string `yaml:"cache_dir"`
}

func init() {

	log.SetOutput(os.Stdout)

	C = &Config{}

	if data, err := ioutil.ReadFile("config.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), C)
	}

	log.Println("加载配置", C)
}
