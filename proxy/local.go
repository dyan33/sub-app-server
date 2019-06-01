package proxy

import (
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"regexp"
)

func LocalProxy(port string) {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(":443$"))).HandleConnect(goproxy.AlwaysMitm)

	//拦截浏览器请求
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

		log.Println("拦截请求 ===>", req.URL.String())

		return req, nil
	})

	log.Println("LocalProxy running on", port)

	if err := http.ListenAndServe(port, proxy); err != nil {
		log.Println(err)
	}

	//adb shell settings put global http_proxy 192.168.50.165:8090
	//adb shell settings delete global http_proxy
	//adb shell settings delete global global_http_proxy_host
	//adb shell settings delete global global_http_proxy_port

}
