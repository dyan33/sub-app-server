package scripts

import (
	"SubAppServer/config"
	"log"
	"os/exec"
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

	script := config.Cfg.Get(s.app.Operator)

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
	deviceid := s.app.Deviceid

	log.Println("执行脚本", exe, name, lang, timezone, proxy, deviceid)

	command := exec.Command(exe, name, lang, timezone, proxy, deviceid)
	command.Dir = script.Dir

	out, err := command.CombinedOutput()

	return string(out), err
}
