package speakall

import (
	"log"
	. "speakall/config"
	"speakall/db"
	"speakall/web"
	"speakall/ws"
)

func Start() error {

	log.Println("############### start DBServer")
	err := db.Listen(Config.Database.Path)
	if err != nil {
		return err
	}

	log.Println("############### start WSServer")
	ws.Listen("/ws/")

	log.Println("############### start HTTPServer")
	port := Config.Web.Port
	staticDir := Config.Web.Root
	web.Listen(staticDir, port)

	return nil
}
