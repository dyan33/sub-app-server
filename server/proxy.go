package server

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"regexp"
	"sub-app-server/config"
	"sync"
	"time"
)

const (
	BEGIN   = "begin"
	NETWORK = "network"
	SMS     = "sms"
)

type AppProxy struct {
	id *ID

	tasks *sync.Map

	port string
	send chan *HttpRequest
	sms  chan string

	proxy *goproxy.ProxyHttpServer
}

func (a *AppProxy) clean() {

	a.tasks.Range(func(key, value interface{}) bool {

		a.tasks.Delete(key)

		close(value.(chan HttpResponse))

		return true
	})
}

func (a *AppProxy) process(req *http.Request) (resp *http.Response) {
	var flag string

	start := time.Now()

	defer func() {
		log.Println(a.port, req.Method, resp.StatusCode, req.URL.String(), time.Now().Sub(start), flag)
	}()

	//缓存加载
	if resp = loadCache(req); resp == nil {

		id, rev := a.id.get(), make(chan HttpResponse)

		a.tasks.Store(id, rev)

		a.send <- makeRequest(id, req)

		if response, ok := <-rev; ok {

			//缓存响应
			if cacheResponse(req, response) {
				flag = "[cached]"
			}

			resp = makeResponse(req, response)

			return
		}

		resp = goproxy.NewResponse(req, "text/plain", 555, "close")
	} else {
		flag = "[load cache]"
	}

	return
}

func (a *AppProxy) doResp(data string) {
	response := HttpResponse{}

	if err := json.Unmarshal([]byte(data), &response); err != nil {
		log.Println("parse json error:", err)
		return
	}

	if recive, ok := a.tasks.Load(response.Id); ok {

		recive.(chan HttpResponse) <- response

		a.tasks.Delete(response.Id)

		close(recive.(chan HttpResponse))
		return
	}

	log.Println(a.port, "not found response", response.Id)
}

func (a *AppProxy) doWork(data string) {
	app := AppInfo{}
	if err := json.Unmarshal([]byte(data), &app); err != nil {
		log.Println("parse json error:", err)
		return
	}

	proxy := fmt.Sprintf("127.0.0.1%s", a.port)

	log.Println(a.port, "start call script！", app)

	info, err := NewBrowerScript(app, proxy).Run()

	log.Println(a.port, "call script end!\n", info, "\n", err)

}

func (a *AppProxy) handle(message []byte, close func()) {

	//防止a.sms没有消耗 关闭chan崩溃
	defer func() { _ = recover() }()

	typ := gjson.GetBytes(message, "type").String()
	data := gjson.GetBytes(message, "data").String()

	switch typ {
	case BEGIN:
		a.doWork(data)

		//脚本执行完毕关闭socket连接
		close()
		break

	case NETWORK:
		a.doResp(data)
		break

	case SMS:
		a.sms <- data
		break

	default:
		log.Println(a.port, "can't handle data!", string(message))

		//不能处理的任务关闭连接
		close()
	}

}

func (a *AppProxy) run() {

	go func() {

		for {
			select {

			case ch := <-ConnChan:

				log.Println(a.port, "start proxy channel!")

				a.sms = make(chan string, 100)
				a.id.set(0)

				ch.Run(a.port, a.send, a.handle)

				close(a.sms)

				a.clean()

				log.Println(a.port, "stop proxy channel!")

			//释放请求
			case _ = <-a.send:

				a.clean()
			}
		}

	}()

	log.Println("ProxyServer running", a.port)
	if err := http.ListenAndServe(a.port, a.proxy); err != nil {
		log.Println(err)
	}
}

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

func newAppProxy(port string) *AppProxy {

	app := &AppProxy{
		tasks: &sync.Map{},
		port:  port,
		send:  make(chan *HttpRequest),
		id:    newID(),
		proxy: goproxy.NewProxyHttpServer(),
	}

	app.proxy.Verbose = false
	app.proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(":443$"))).HandleConnect(goproxy.AlwaysMitm)
	app.proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (request *http.Request, response *http.Response) {

		//发送短信息
		if "sms" == req.Host {

			if app.sms != nil {
				if text, ok := <-app.sms; ok {
					return nil, goproxy.NewResponse(req, "text/plain", 200, text)
				}
			}
			return nil, goproxy.NewResponse(req, "text/plain", 555, "")

		}

		//转发请求
		if isForward(req) {

			return nil, app.process(req)

		}
		return nil, goproxy.NewResponse(req, "text/plain", 404, "")
	})

	return app
}

func RunProxy(port string) {
	newAppProxy(port).run()
}
