package task

import "github.com/gorilla/websocket"

var TasksChan chan Task

type Task struct {
	Url  string
	Html string
	Conn *websocket.Conn
}

func init() {
	TasksChan = make(chan Task, 4)
}
