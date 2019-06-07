package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

var ConnChan = make(chan *Channel)

/*
客户端连接到服务端
*/
func wsHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}

	ConnChan <- NewChannel(conn)

}

func WebServer(port string) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/ws", wsHandler)

	r.GET("/", func(c *gin.Context) {

		c.HTML(200, "index.html", nil)
	})

	log.Println("WebServer running", port)

	if err := r.Run(port); err != nil {
		log.Println(err)
	}

}
