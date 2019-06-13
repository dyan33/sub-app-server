package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
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

type LogWriter struct {
	m       *sync.Mutex
	logname string
	logfile io.Writer
}

func init() {

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(mw)

	C = &Config{}

	if data, err := ioutil.ReadFile("config.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), C)
	}

	log.Println("加载配置", C)
}
