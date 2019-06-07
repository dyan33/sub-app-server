package server

import (
	"github.com/gorilla/websocket"
	"log"
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

type Channel struct {
	conn *websocket.Conn
	send chan *HttpRequest

	stop chan int

	handle func(data []byte, close func())

	wait *sync.WaitGroup

	name string
}

func (c *Channel) reader() {

	defer func() {
		_ = c.conn.Close()
		close(c.stop)
		c.wait.Done()
		log.Println(c.name, "close reader !")
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf(c.name, "read message error: %v", err)
			}
			break
		}

		go c.handle(data, c.close)

	}

}

func (c *Channel) writer() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		_ = c.conn.Close()
		ticker.Stop()
		c.wait.Done()

		log.Println(c.name, "close writer !")
	}()

	for {
		select {

		case req := <-c.send:

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteJSON(req); err != nil {
				log.Println(c.name, "socket write error", err)
			}

		case <-c.stop:
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}

}

func (c *Channel) close() {

	_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
	_ = c.conn.Close()
}

func (c *Channel) Run(name string, send chan *HttpRequest, handle func(data []byte, close func())) {

	c.name = name
	c.send = send
	c.handle = handle

	go c.reader()
	go c.writer()

	c.wait.Wait()
	log.Println(c.name, "close channel !")
}

func NewChannel(conn *websocket.Conn) *Channel {

	ch := &Channel{
		conn: conn,
		stop: make(chan int),
		wait: &sync.WaitGroup{},
	}

	ch.wait.Add(2)

	return ch

}
