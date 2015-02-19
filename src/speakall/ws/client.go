package ws

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type client struct {
	Id string
	ws *websocket.Conn
}

func NewClient(ws *websocket.Conn) *client {
	return &client{
		Id: uuid.NewV4().String(),
		ws: ws,
	}
}

func (c *client) start(msgCh chan *message, removeCh chan *client) {
	c.send(createOpenMessage(c.Id))
	for {
		msg := &message{}
		err := websocket.JSON.Receive(c.ws, msg)
		if err == nil {
			msg.ClientId = c.Id
			msg.Date = time.Now().String()
			msgCh <- msg
		} else {
			removeCh <- c
			return
		}
		c.send(createOpenMessage(c.Id))
	}
}

func (c *client) send(msg *message) {
	err := websocket.JSON.Send(c.ws, msg)
	if err != nil {
		log.Println(err)
	}
}
