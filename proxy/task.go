package proxy

import (
	"net/http"
	"sub-app-server/config"
)

func isForward(req *http.Request) bool {

	scheme := req.URL.Scheme

	if scheme == "http" || scheme == "https" {

		host := req.Host

		for _, h := range config.C.Igonre {
			if h == host {
				return false
			}
		}

		return true
	}

	return false

}

//func TaskProxy(port string) {
//
//	proxy := goproxy.NewProxyHttpServer()
//	proxy.Verbose = false
//	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(":443$"))).HandleConnect(goproxy.AlwaysMitm)
//
//	app := server.NewSocketClient(port)
//
//	app.Run()
//
//	//拦截浏览器请求
//	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {
//
//		if "sms" == req.Host {
//
//			if text, ok := <-app.Sms(); ok {
//				return nil, goproxy.NewResponse(req, "text/plain", 200, text)
//			}
//			return nil, goproxy.NewResponse(req, "text/plain", 555, "")
//		}
//
//		if isForward(req) {
//			return nil, app.Process(req)
//		}
//
//		return nil, goproxy.NewResponse(req, "text/plain", 404, "")
//
//	})
//
//	log.Println("start TaskProxy", port)
//
//	if err := http.ListenAndServe(port, proxy); err != nil {
//		log.Println(err)
//	}
//}
