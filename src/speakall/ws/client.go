package ws

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"time"
)

type client struct {
	Id string
	ws *websocket.Conn
}

type message struct {
	Content  string
	Category string
	Date     string
}

func NewClient(ws *websocket.Conn) *client {
	return &client{
		Id: uuid.NewV4().String(),
		ws: ws,
	}
}

func (c *client) start(msgCh chan *message, removeCh chan *client) {
	for {
		msg := &message{}
		err := websocket.JSON.Receive(c.ws, msg)
		if err == nil {
			msg.Date = time.Now().String()
			msgCh <- msg
		} else {
			removeCh <- c
			return
		}
	}
}

func (c *client) send(msg *message) {
	err := websocket.JSON.Send(c.ws, msg)
	if err != nil {
		fmt.Println(err)
	}
}
