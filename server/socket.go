package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HttpRequest struct {
	Id     int64             `json:"id"`
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   []byte            `json:"body"`
}

type HttpResponse struct {
	Id      int64       `json:"id"`
	Code    int         `json:"code"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

//type SocketClient struct {
//	port string
//	name string
//
//	id    *ID
//	tasks *sync.Map
//
//	send chan *HttpRequest
//	sms  chan string
//
//	conn *websocket.Conn
//}

func makeRequest(id int64, req *http.Request) *HttpRequest {

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

////清空tasks
//func (w *SocketClient) cleanTask() {
//
//	log.Println(w.name, "clean tasks!")
//
//	w.tasks.Range(func(key, value interface{}) bool {
//
//		w.tasks.Delete(key)
//
//		close(value.(chan HttpResponse))
//
//		return true
//	})
//}
//
////写
//func (w *SocketClient) wirteSocket(stop <-chan int) {
//
//	ticker := time.NewTicker(pingPeriod)
//
//	defer func() {
//
//		ticker.Stop()
//
//		_ = w.conn.Close()
//
//		log.Println(w.name, "close write socket!")
//	}()
//
//	for {
//		select {
//
//		case reqeust := <-w.send:
//
//			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
//
//			if err := w.conn.WriteJSON(reqeust); err != nil {
//				log.Println("socket write error", err)
//			}
//
//		case <-stop:
//			_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
//			return
//
//		case <-ticker.C:
//			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//
//	}
//}
//
////读
//func (w *SocketClient) readSocket(stop chan int) {
//
//	defer func() {
//
//		close(stop)
//
//		_ = w.conn.Close()
//
//		log.Println(w.name, "close read socket!")
//
//	}()
//
//	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
//	w.conn.SetPongHandler(func(string) error { _ = w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
//
//	for {
//
//		_, message, err := w.conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("read message error: %v", err)
//			}
//			break
//		}
//
//		mType := gjson.GetBytes(message, "type").String()
//		data := []byte(gjson.GetBytes(message, "data").String())
//
//		switch mType {
//		case "begin":
//			go w.callScript(data)
//			break
//		case "network":
//			go w.doResponse(data)
//			break
//		case "sms":
//			go func(data []byte) { defer func() { recover() }(); w.sms <- string(data) }(data)
//			break
//
//		default:
//			log.Println(w.name, "not handle data!", string(message))
//			return
//		}
//
//	}
//}
//
////响应
//func (w *SocketClient) doResponse(data []byte) {
//
//	response := HttpResponse{}
//
//	if err := json.Unmarshal(data, &response); err != nil {
//		log.Println("parse json error:", err)
//		return
//	}
//
//	if recive, ok := w.tasks.Load(response.Id); ok {
//
//		recive.(chan HttpResponse) <- response
//
//		w.tasks.Delete(response.Id)
//
//		close(recive.(chan HttpResponse))
//		return
//	}
//
//	fmt.Println("not found response", response.Id)
//}
//
////执行脚本
//func (w *SocketClient) callScript(data []byte) {
//
//	app := AppInfo{}
//	if err := json.Unmarshal(data, &app); err != nil {
//		log.Println("parse json error:", err)
//		return
//	}
//
//	proxy := fmt.Sprintf("127.0.0.1%s", w.port)
//
//	log.Println(w.name, "start call script！", app)
//
//	info, err := NewBrowerScript(app, proxy).Run()
//
//	log.Println(w.name, "call script end!\n", info, "\n", err)
//
//	_ = w.conn.WriteMessage(websocket.CloseMessage, nil)
//	_ = w.conn.Close()
//}
//
//func (w *SocketClient) Run() {
//
//	go func() {
//
//		for {
//
//			select {
//
//			case w.conn = <-SocketChan:
//
//				//初始化
//				w.id.set(0)
//				w.sms = make(chan string, 100) //缓存100条通知
//
//				log.Println(w.name, "start task!")
//
//				stopChan := make(chan int)
//
//				go w.wirteSocket(stopChan)
//				go w.readSocket(stopChan)
//
//				_ = <-stopChan
//
//				close(w.sms)
//
//				w.cleanTask()
//
//				log.Println(w.name, "over task!")
//
//			case _ = <-w.send:
//				w.cleanTask()
//			}
//		}
//	}()
//}
//
////处理请求转发
//func (w *SocketClient) Process(req *http.Request) (resp *http.Response) {
//
//	var flag string
//
//	start := time.Now()
//
//	defer func() {
//		log.Println(w.name, req.Method, resp.StatusCode, req.URL.String(), time.Now().Sub(start), flag)
//	}()
//
//	//缓存加载
//	if resp = loadCache(req); resp == nil {
//
//		id, rev := w.id.get(), make(chan HttpResponse)
//
//		w.tasks.Store(id, rev)
//
//		w.send <- makeRequest(id, req)
//
//		if response, ok := <-rev; ok {
//
//			//缓存响应
//			if cacheResponse(req, response) {
//				flag = "[cached]"
//			}
//
//			resp = makeResponse(req, response)
//
//			return
//		}
//
//		resp = goproxy.NewResponse(req, "text/plain", 555, "close")
//	} else {
//		flag = "[load cache]"
//	}
//
//	return
//
//}
//
////短信
//func (w *SocketClient) Sms() chan string {
//	return w.sms
//}
//
//func NewSocketClient(port string) *SocketClient {
//
//	return &SocketClient{
//		port:  port,
//		name:  "[proxy" + port + "]",
//		tasks: &sync.Map{},
//		send:  make(chan *HttpRequest),
//		id: &ID{
//			mutex: &sync.Mutex{},
//		},
//	}
//}
