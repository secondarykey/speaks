package ws

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type Server struct {
	clients  map[string]*client
	addCh    chan *client
	removeCh chan *client
	msgCh    chan *message
}

func Listen(path string) {
	server := &Server{
		clients:  make(map[string]*client),
		addCh:    make(chan *client),
		removeCh: make(chan *client),
		msgCh:    make(chan *message),
	}
	go server.listen(path)
}

func (s *Server) add(c *client) {
	s.clients[c.Id] = c
}

func (s *Server) remove(c *client) {
	delete(s.clients, c.Id)
}

func (s *Server) sendMessage(msg *message) {

	for _, c := range s.clients {
		client := c
		if msg.Category == client.Category {
			go func() {
				client.send(msg)
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

func (s *Server) listen(pattern string) {
	http.Handle(pattern, s.websocketHandler())
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
}
