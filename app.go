package main

import (
	"fmt"
	"sub-app-server/config"
	"sub-app-server/proxy"
	"sub-app-server/server"
)

func run(server func(string), ports []int) bool {

	b := false

	for _, port := range ports {
		go server(fmt.Sprintf(":%d", port))
		b = true
	}
	return b

}

func main() {

	//代理
	run(proxy.TaskProxy, config.C.Proxy)

	//web
	run(server.WebServer, config.C.Server)

	proxy.LocalProxy(":8030")

	<-make(chan struct{})

}
