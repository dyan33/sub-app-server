package server

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Channel struct {
	conn *websocket.Conn
	send chan *HttpRequest

	stop chan int

	handle func(data []byte, close func())

	wait *sync.WaitGroup
}

func (c *Channel) reader() {

	defer func() {
		_ = c.conn.Close()
		close(c.stop)
		c.wait.Done()
		log.Println("close reader !")
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message error: %v", err)
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

		log.Println("close writer !")
	}()

	for {
		select {

		case req := <-c.send:

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteJSON(req); err != nil {
				log.Println("socket write error", err)
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

func (c *Channel) Run() {

	go c.reader()
	go c.writer()

	c.wait.Wait()
	log.Println("close channel !")
}

func NewChannel(conn *websocket.Conn, send chan *HttpRequest, handle func(data []byte, close func())) *Channel {

	ch := &Channel{
		conn,
		send,
		make(chan int),
		handle,
		&sync.WaitGroup{},
	}

	ch.wait.Add(2)

	return ch

}
