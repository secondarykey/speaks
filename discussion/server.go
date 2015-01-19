package discussion

import (
	"../db"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Server struct {
	clients  map[string]*client
	addCh    chan *client
	removeCh chan *client
	msgCh    chan *message
}

/**
 * Create Server
 **/
func NewServer() *Server {
	return &Server{
		clients:  make(map[string]*client),
		addCh:    make(chan *client),
		removeCh: make(chan *client),
		msgCh:    make(chan *message),
	}
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
		go func() {
			client.send(msg)
		}()
	}

	db.Select()
}

func (s *Server) WebsocketHandler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		c := NewClient(ws)
		s.addCh <- c
		log.Println("Client Start[" + c.Id + "]")
		c.start(s.msgCh, s.removeCh)
	})
}

func (s *Server) Listen(pattern string) {
	http.Handle(pattern, s.WebsocketHandler())
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
