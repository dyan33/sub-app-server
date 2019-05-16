package main

import (
	"SubAppServer/config"
	"SubAppServer/proxy"
	"SubAppServer/server"
	"fmt"
)

func main() {

	for i := 0; i < config.Cfg.ProxyNum; i++ {
		go proxy.OrangeProxy(fmt.Sprintf(":%d", config.Cfg.ProxyPort+i))
	}

	server.WebServer(fmt.Sprintf(":%d", config.Cfg.WebPort))

}
