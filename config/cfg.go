package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var Cfg *Config

type AppInfo struct {
	Operator string `json:"operator"`
	Deviceid string `json:"deviceid"`
	TimeZone string `json:"timezone"`
	Lang     string `json:"lang"`
}

type Script struct {
	Exe  string `yaml:"exe"`
	Dir  string `yaml:"dir"`
	Name string `yaml:"name"`
}

type Config struct {
	Server  []int             `yaml:"server"`
	Proxy   []int             `yaml:"proxy"`
	Hosts   []string          `yaml:"hosts"`
	Scripts map[string]Script `yaml:"scripts"`

	mutex sync.Mutex
}

func (c Config) String() string {

	return fmt.Sprintf(`

	server: %d
	 proxy: %d 
     hosts: %s
   scripts: %s
`, c.Server, c.Proxy, c.Hosts, c.Scripts)
}

func (c Config) Get(operator string) Script {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.Scripts[operator]
}

func init() {

	log.SetOutput(os.Stdout)

	Cfg = &Config{}

	if data, err := ioutil.ReadFile("config.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), Cfg)
	}

	Cfg.mutex = sync.Mutex{}

	log.Println("加载配置", Cfg)
}
