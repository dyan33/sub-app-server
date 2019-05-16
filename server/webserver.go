package server

import (
	"SubAppServer/task"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Duration(time.Second * 5),
}

/*
客户端连接到服务端
*/
func wsHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}

	for {
		msgType, msgData, err := conn.ReadMessage()

		if err != nil {
			log.Println("cant read message:", err)

			switch err.(type) {
			case *websocket.CloseError:
				return
			default:
				continue
			}
		}

		if msgType == websocket.TextMessage {

			text := string(msgData)

			location := gjson.Get(text, "location").String()
			html := gjson.Get(text, "html").String()
			vid := gjson.Get(text, "vid").String()

			task.TasksChan <- task.Task{
				Url:  location,
				Html: html,
				Conn: conn,
			}

			log.Println("接收任务", vid, location)
		}

	}

}

func WebServer(port string) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/ws", wsHandler)

	log.Println("WebServer running on", port)

	if err := r.Run(port); err != nil {
		log.Println(err)
	}

}
