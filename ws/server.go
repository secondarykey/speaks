package ws

import (
	"net/http"

	"github.com/secondarykey/speaks/db"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients  map[string]*client
	username map[int]string
	addCh    chan *client
	removeCh chan *client
	msgCh    chan *message
}

func Listen(path string) error {

	server := &Server{
		clients:  make(map[string]*client),
		username: make(map[int]string),
		addCh:    make(chan *client),
		removeCh: make(chan *client),
		msgCh:    make(chan *message),
	}
	return server.listen(path)
}

func (s *Server) add(c *client) {
	s.clients[c.Id] = c
	s.sendMessage(createAddUserMessage(c.UserId, c.Id, s.getUserName(c.UserId)))
}

func (s *Server) remove(c *client) {
	delete(s.clients, c.Id)
	s.sendMessage(createDeleteUserMessage(c.Id))
}

func (s *Server) sendMessage(msg *message) {

	msg.UserName = s.getUserName(msg.UserId)

	for _, c := range s.clients {
		client := c
		if msg.Project == client.Project && msg.Category == client.Category {
			go func() {
				client.send(msg)
			}()
		} else {
			if msg.Type == "Message" {
				go func() {
					client.send(createBadgeMessage(msg.Project, msg.Category))
				}()
			} else if msg.Type == "AddUser" || msg.Type == "DeleteUser" {
				go func() {
					client.send(msg)
				}()
			}
		}
	}
}

func (s *Server) websocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		c := newClient(ws)
		for _, oc := range s.clients {
			c.send(createAddUserMessage(oc.UserId, oc.Id, s.getUserName(oc.UserId)))
		}
		s.addCh <- c
		c.start(s.msgCh, s.removeCh)
	})
}

func (s *Server) listen(pattern string) error {

	http.Handle(pattern, s.websocketHandler())
	//goroutine
	go func() {
		for {
			select {
			case c := <-s.addCh:
				s.add(c)
			case c := <-s.removeCh:
				s.remove(c)
			case m := <-s.msgCh:
				s.sendMessage(m)
			}
		}
	}()

	return nil
}

func (s *Server) getUserName(userId int) string {
	name := "someone"
	if s.username[userId] == "" {
		wk, err := db.GetUserName(userId)
		if err == nil {
			name = wk
			s.username[userId] = name
		}
	} else {
		name = s.username[userId]
	}
	return name
}
