package speakall

import (
	"log"
	. "speakall/config"
	"speakall/db"
	"speakall/web"
	"speakall/ws"
)

func Listen() {

	log.Println("############### start DBServer")
	db.Listen(Config.Database.Path)

	log.Println("############### start WSServer")
	ws.Listen("/ws/")

	log.Println("############### start HTTPServer")

	port := Config.Web.Port
	staticDir := Config.Web.Root
	web.Listen(staticDir, port)
}
