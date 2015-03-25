package ws

import (
	"golang.org/x/net/websocket"
	"net/http"
	"speakall/db"
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
}

func (s *Server) remove(c *client) {
	delete(s.clients, c.Id)
}

func (s *Server) sendMessage(msg *message) {

	name := "someone"
	if s.username[msg.UserId] == "" {
		wk, err := db.GetUserName(msg.UserId)
		if err == nil {
			name = wk
			s.username[msg.UserId] = name
		}
	} else {
		name = s.username[msg.UserId]
	}
	msg.UserName = name

	for _, c := range s.clients {
		client := c
		if msg.Category == client.Category {
			go func() {
				client.send(msg)
			}()
		} else {
			//notify badge
			go func() {
				client.send(createBadgeMessage(msg.Category))
			}()
		}
	}
}

func (s *Server) websocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		c := newClient(ws)
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
