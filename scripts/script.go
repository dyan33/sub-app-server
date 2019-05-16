package scripts

import (
	"SubAppServer/config"
	"os/exec"
)

type orange struct {
	port string
	exe  string
}

func NewOrange(port string) *orange {

	return &orange{
		port: port,
		exe:  config.Cfg.Python,
	}
}

func (o *orange) Run(url string) (string, error) {

	out, err := exec.Command(
		o.exe,
		config.Cfg.OrangeScript,
		url,
		o.port).CombinedOutput()

	return string(out), err

}
