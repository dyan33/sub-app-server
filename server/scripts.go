package server

import (
	"os/exec"
	"sub-app-server/config"
)

type BrowerScript struct {
	app   config.AppInfo
	proxy string
}

func NewBrowerScript(app config.AppInfo, proxy string) *BrowerScript {

	return &BrowerScript{
		app,
		proxy,
	}
}

func (s *BrowerScript) Run() (string, error) {

	script := config.Cfg.Get(s.app.OperatorCode)

	exe := script.Exe
	//脚本名字
	name := script.Name

	//语言
	lang := s.app.Lang

	//时区
	timezone := s.app.TimeZone

	//代理
	proxy := s.proxy

	//设备id
	android := s.app.AndroidId

	command := exec.Command(exe, name, lang, timezone, proxy, android)
	command.Dir = script.Dir

	out, err := command.CombinedOutput()

	return string(out), err
}
