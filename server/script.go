package server

import (
	"os/exec"
	"sub-app-server/config"
)

type AppInfo struct {
	AndroidId    string `json:"android_id"`
	Version      string `json:"version"`
	SdkVersion   string `json:"sdk_version"`
	DeviceName   string `json:"device_name"`
	OperatorName string `json:"operator_name"`
	OperatorCode string `json:"operator_code"`
	PackageName  string `json:"package_name"`
	Network      string `json:"network"`

	TimeZone string `json:"timezone"`
	Lang     string `json:"lang"`
}

type BrowerScript struct {
	app   AppInfo
	proxy string
}

func NewBrowerScript(app AppInfo, proxy string) *BrowerScript {

	return &BrowerScript{
		app,
		proxy,
	}
}

func (s *BrowerScript) Run() (string, error) {

	script := config.S.Get(s.app.OperatorCode)

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
