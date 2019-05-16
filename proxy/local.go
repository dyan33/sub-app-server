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

}
