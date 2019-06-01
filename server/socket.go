package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/elazarl/goproxy.v1"
	"io/ioutil"
	"log"
	"net/http"
	"sub-app-server/config"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var SocketChan chan *websocket.Conn

func init() {
	SocketChan = make(chan *websocket.Conn, 4)
}

type HttpRequest struct {
	Id     int               `json:"id"`
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   []byte            `json:"body"`
}

type HttpResponse struct {
	Id      int         `json:"id"`
	Code    int         `json:"code"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

type SocketClient struct {
	port  string
	name  string
	tasks sync.Map

	send chan *HttpRequest

	conn *websocket.Conn
}

func makeRequest(id int, req *http.Request) *HttpRequest {

	body, _ := ioutil.ReadAll(req.Body)

	request := &HttpRequest{
		Id:     id,
		Method: req.Method,
		Header: map[string]string{},
		Url:    req.URL.String(),
		Body:   body,
	}

	for key, value := range req.Header {

		if key == "Proxy-Connection" {
			continue
		}

		request.Header[key] = value[0]
	}

	return request
}

func makeResponse(r *http.Request, response HttpResponse) *http.Response {

	buf := bytes.NewBuffer(response.Body)

	return &http.Response{
		Request:          r,
		TransferEncoding: r.TransferEncoding,
		Header:           response.Headers,
		StatusCode:       response.Code,
		ContentLength:    int64(buf.Len()),
		Body:             ioutil.NopCloser(buf),
	}

}

//清空tasks
func (w *SocketClient) cleanTask() {
	w.tasks.Range(func(key, value interface{}) bool {

		w.tasks.Delete(key)

		close(value.(chan HttpResponse))

		return true
	})
}

//写
func (w *SocketClient) wirteSocket(stop <-chan int) {

	ticker := time.NewTicker(pingPeriod)

	defer func() {

		ticker.Stop()

		_ = w.conn.Close()

		log.Println(w.name, "关闭写入!")
	}()

	for {
		select {

		case reqeust := <-w.send:

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := w.conn.WriteJSON(reqeust); err != nil {
				log.Println("socket write error", err)
			}

		case <-stop:
			_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return

		case <-ticker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

//读
func (w *SocketClient) readSocket(stop chan int) {

	defer func() {

		close(stop)

		_ = w.conn.Close()

		log.Println(w.name, "关闭读取!")

	}()

	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error { _ = w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		_, message, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message error: %v", err)
			}
			break
		}

		mType := gjson.GetBytes(message, "type").String()
		data := []byte(gjson.GetBytes(message, "data").String())

		switch mType {
		case "start":
			go w.runScript(data)
			break
		case "running":
			w.doResponse(data)
			break
		default:
			log.Println(w.name, "不能处理数据!", string(message))
		}

	}
}

//响应
func (w *SocketClient) doResponse(data []byte) {

	response := HttpResponse{}

	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("解析Json错误:", err)
		return
	}

	if recive, ok := w.tasks.Load(response.Id); ok {

		recive.(chan HttpResponse) <- response

		w.tasks.Delete(response.Id)

		close(recive.(chan HttpResponse))
	}
}

//执行脚本
func (w *SocketClient) runScript(data []byte) {

	app := config.AppInfo{}
	if err := json.Unmarshal(data, &app); err != nil {
		log.Println("解析Json错误:", err)
		return
	}

	proxy := fmt.Sprintf("127.0.0.1%s", w.port)

	log.Println(w.name, "开始执行脚本！")

	info, err := NewBrowerScript(app, proxy).Run()

	log.Println(w.name, "脚本已完成 ===>", info, err)

	_ = w.conn.WriteMessage(websocket.CloseMessage, nil)
	_ = w.conn.Close()
}

func (w *SocketClient) Run() {

	go func() {

		for {

			select {

			case w.conn = <-SocketChan:

				log.Println(w.name, "开始任务!")

				stopChan := make(chan int)

				go w.wirteSocket(stopChan)
				go w.readSocket(stopChan)

				_ = <-stopChan

				w.cleanTask()

				log.Println(w.name, "结束任务!")

			case _ = <-w.send:
				w.cleanTask()
			}
		}
	}()
}

//处理请求转发
func (w *SocketClient) Process(req *http.Request) (resp *http.Response) {

	start, id := time.Now(), time.Now().Nanosecond()

	defer func() { log.Println(w.name, req.Method, resp.StatusCode, req.URL.String(), time.Now().Sub(start)) }()

	rev := make(chan HttpResponse)
	w.tasks.Store(id, rev)

	w.send <- makeRequest(id, req)

	if response, ok := <-rev; ok {

		resp = makeResponse(req, response)

		return
	}

	resp = goproxy.NewResponse(req, "text/plain", 555, "close")

	return

}

func NewSocketClient(port string) *SocketClient {

	return &SocketClient{
		port:  port,
		name:  "[代理" + port + "]",
		tasks: sync.Map{},
		send:  make(chan *HttpRequest),
	}
}
