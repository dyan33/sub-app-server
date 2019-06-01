package proxy

import (
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"regexp"
	"sub-app-server/config"
	"sub-app-server/server"
)

func isForward(req *http.Request) bool {

	scheme := req.URL.Scheme

	if scheme == "http" || scheme == "https" {

		host := req.Host

		for _, h := range config.Cfg.Hosts {
			if h == host {
				return false
			}
		}

		return true
	}

	return false

}

func TaskProxy(port string) {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(":443$"))).HandleConnect(goproxy.AlwaysMitm)

	app := server.NewSocketClient(port)

	app.Run()

	//拦截浏览器请求
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

		if isForward(req) {
			return nil, app.Process(req)
		}

		return nil, goproxy.NewResponse(req, "text/plain", 404, "")

	})

	log.Println("start TaskProxy", port)

	if err := http.ListenAndServe(port, proxy); err != nil {
		log.Println(err)
	}
}
