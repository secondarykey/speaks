package ws

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/secondarykey/speaks/db"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

type client struct {
	Id       string
	UserId   int
	Project  string
	Category string
	ws       *websocket.Conn
}

func newClient(ws *websocket.Conn) *client {

	url := ws.Request().URL
	paths := strings.Split(url.String(), "/")
	userId, _ := strconv.Atoi(paths[2])

	return &client{
		Id:       uuid.NewV4().String(),
		UserId:   userId,
		Project:  "Speaks",
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

			log.Println(msg.Type)

			if msg.Type == "Delete" {
				msgCh <- msg
			} else if msg.Type != "Change" {
				t := time.Now()
				msg.Created = t.Format("2006/01/02 15:04:05")
				go db.InsertMessage(msg.UserId, msg.Project, msg.Category, msg.Content, msg.Created)
				msgCh <- msg
			} else {
				c.Project = msg.Project
				c.Category = msg.Category
			}
		} else {

			log.Println(err)

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
