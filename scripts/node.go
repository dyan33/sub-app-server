package scripts

import (
	"SubAppServer/config"
	"log"
	"os/exec"
)

type BrowerScript struct {
	config.Script
	proxy string
}

func NewBrowerScript(script config.Script, proxy string) *BrowerScript {

	return &BrowerScript{
		script,
		proxy,
	}
}

func (s *BrowerScript) Run() (string, error) {

	log.Println("执行脚本", s.Script.Exe, s.Script.Name, s.proxy)

	command := exec.Command(s.Script.Exe, s.Script.Name, s.proxy)
	command.Dir = s.Script.Dir

	out, err := command.CombinedOutput()

	return string(out), err
}
