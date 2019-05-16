package proxy

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"net/url"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", basicAuth(username, password)))
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func SimtestProxy(port string) {
	username, password := "mauritius", "Ux5vW5qw"

	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Verbose = false

	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		SetBasicAuth(username, password, req)

		return url.Parse("http://91.220.77.154:8090")
	}

	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler("http://91.220.77.154:8090", func(req *http.Request) {
		SetBasicAuth(username, password, req)
	})

	middleProxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		SetBasicAuth(username, password, req)
		return req, nil
	}))

	log.Println("SimTest running on", port)

	if err := http.ListenAndServe(port, middleProxy); err != nil {
		log.Println(err)
	}
}
