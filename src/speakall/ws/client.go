package ws

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"speakall/db"
	"time"
)

type client struct {
	Id       string
	Category string
	ws       *websocket.Conn
}

func newClient(ws *websocket.Conn) *client {
	return &client{
		Id:       uuid.NewV4().String(),
		Category: "Dashboard",
		ws:       ws,
	}
}

func (c *client) start(msgCh chan *message, removeCh chan *client) {
	c.send(createOpenMessage(c.Id))
	for {
		msg := &message{}
		err := websocket.JSON.Receive(c.ws, msg)
		if err == nil {
			if msg.Type != "Change" {
				t := time.Now()
				msg.Created = t.Format("2006/01/02 15:04:05")
				go db.InsertMessage(msg.UserId, msg.Category, msg.Content, msg.Created)
				msgCh <- msg
			} else {
				c.Category = msg.Category
			}
		} else {
			removeCh <- c
			return
		}
	}
}

func (c *client) send(msg *message) {
	err := websocket.JSON.Send(c.ws, msg)
	if err != nil {
	}
}
