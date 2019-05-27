package proxy

import (
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
)

func LocalProxy(port string) {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	log.Println("LocalProxy running on", port)

	if err := http.ListenAndServe(port, proxy); err != nil {
		log.Println(err)
	}

	//adb shell settings put global http_proxy 192.168.50.165:8090
	//adb shell settings delete global http_proxy
	//adb shell settings delete global global_http_proxy_host
	//adb shell settings delete global global_http_proxy_port

}
