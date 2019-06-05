package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type Script struct {
	Exe  string `yaml:"exe"`
	Dir  string `yaml:"dir"`
	Name string `yaml:"name"`
}

type ScriptConfig struct {
	Scripts map[string]Script `yaml:"scripts"`
	mutex   *sync.Mutex
}

func (s *ScriptConfig) Get(key string) *Script {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if value, ok := s.Scripts[key]; ok {
		return &value
	}
	return nil

}

var S *ScriptConfig

func init() {
	S = &ScriptConfig{}

	if data, err := ioutil.ReadFile("script.yaml"); err == nil {

		_ = yaml.Unmarshal([]byte(data), S)
	}

	S.mutex = &sync.Mutex{}

	log.Println("加载脚本", S)

}
