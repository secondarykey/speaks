package speaks

import (
	"log"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
	"github.com/secondarykey/speaks/web"
	"github.com/secondarykey/speaks/ws"
)

func init() {
}

var done chan error

func Listen() error {

	var iniFile string

	flag.Parse()
	args := flag.Args()
	leng := len(args)

	switch leng {
	case 0:
		iniFile = "speaks.ini"
	case 1:
		iniFile = args[0]
	default:
		log.Println("Error: too many argument.")
		return
	}

	err := Load(iniFile)
	if err != nil {
		log.Println(err.Error())
		log.Println("Error: Loading initialize file.[" + iniFile + "]")
		return
	}

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
