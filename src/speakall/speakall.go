package speakall

import (
	"log"
	. "speakall/config"
	"speakall/db"
	"speakall/web"
	"speakall/ws"
)

var done chan error

func Start() error {

	log.Println("######## start DBServer")
	path := Config.Database.Path
	ver := Config.Database.Version

	err := db.Listen(path, ver)
	if err != nil {
		return err
	}

	log.Println("######## start WSServer")
	err = ws.Listen("/ws/")
	if err != nil {
		return err
	}

	log.Println("######## start HTTPServer")
	port := Config.Web.Port
	staticDir := Config.Web.Root

	return web.Listen(staticDir, port)
}
