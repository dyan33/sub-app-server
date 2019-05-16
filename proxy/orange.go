package proxy

import (
	"SubAppServer/scripts"
	"SubAppServer/task"
	"github.com/gorilla/websocket"
	"gopkg.in/elazarl/goproxy.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Request struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Form    map[string]string `json:"form"`
}

func getHeader(headers http.Header) map[string]string {

	header := map[string]string{}

	for k, v := range headers {
		header[k] = v[0]
	}
	return header
}

func getFormBody(body io.ReadCloser) map[string]string {

	if data, err := ioutil.ReadAll(body); err == nil {

		if values, err := url.ParseQuery(string(data)); err == nil {

			form := map[string]string{}
			for k, v := range values {
				form[k] = v[0]
			}
			return form
		}
	}
	return nil

}

func OrangeProxy(port string) {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(":443$"))).HandleConnect(goproxy.AlwaysMitm)

	htmlMap := sync.Map{}

	script := scripts.NewOrange(port)

	var conn *websocket.Conn

	//任务获取
	go func() {

		for {
			select {
			case t := <-task.TasksChan:

				htmlMap.Store(t.Url, t)

				conn = t.Conn

				log.Println("执行任务", port, t.Url)

				//执行任务
				out, err := script.Run(t.Url)

				log.Println("PROXY", port, "script", out, err)

				time.Sleep(3 * time.Second)

				log.Println("执行完毕", port, t.Url)

			}
		}

	}()

	tasks := make(chan Request, 4)

	go func() {

		for {

			select {

			case t := <-tasks:

				if conn != nil {
					if err := conn.WriteJSON(t); err != nil {
						log.Println("write error", err)
					}
				}

			}

		}

	}()

	//orange
	{

		//jquery
		proxy.OnRequest(goproxy.DstHostIs("ajax.googleapis.com:443")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

			body, _ := ioutil.ReadFile("static/orange/js/jquery.js")

			return nil, goproxy.NewResponse(req, "application/javascript", 200, string(body))

		})

		//参数
		proxy.OnRequest(goproxy.DstHostIs("notify.dcbprotect.com")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

			r := goproxy.NewResponse(req, "text/plain", 200, "ok")

			if req.Method == "POST" {

				tasks <- Request{
					Url:     req.RequestURI,
					Method:  "POST",
					Headers: getHeader(req.Header),
					Form:    getFormBody(req.Body),
				}

			}

			r.Header.Add("Access-Control-Allow-Origin", "*")

			log.Println("PROXY", port, req.Method, req.RequestURI)

			return req, r

		})

		proxy.OnRequest(goproxy.DstHostIs("enabler.dvbs.com")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

			uri := req.RequestURI

			if strings.HasPrefix(uri, "http://enabler.dvbs.com/session/jscb") { //页面加载请求 GET
				tasks <- Request{
					Method:  req.Method,
					Url:     uri,
					Headers: getHeader(req.Header),
				}
			} else if strings.HasPrefix(uri, "http://enabler.dvbs.com/session/cardpic") { //图片

				body, _ := ioutil.ReadFile("static/orange/image/lp_img.jpeg")

				return nil, NewResponse(req, "image/jpeg", 200, body)

			} else if strings.HasPrefix(uri, "http://enabler.dvbs.com/card/confirmm") { //订阅请求 POST
				tasks <- Request{
					Method:  req.Method,
					Url:     uri,
					Headers: getHeader(req.Header),
					Form:    getFormBody(req.Body),
				}
			}

			log.Println("PROXY", port, req.Method, uri)

			//Html返回
			if value, ok := htmlMap.Load(uri); ok {

				t := value.(task.Task)

				return req, goproxy.NewResponse(req, goproxy.ContentTypeHtml, 200, t.Html)
			}

			return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, 200, "ok")

		})

	}

	log.Println("OrangeProxy running on", port)

	if err := http.ListenAndServe(port, proxy); err != nil {
		log.Println(err)
	}
}
