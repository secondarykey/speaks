package speaks

import (
	"flag"
	"log"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
	"github.com/secondarykey/speaks/web"
	"github.com/secondarykey/speaks/ws"
)

func init() {

	log.SetFlag(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

}

func Usage() {

	// -lv log : Level debug,fatal,error,

	// args[0] sub command : init start version help

	// init -> create .speaks directory
	//    [Database]
	//    version = 0.5
	//    path = data/speaks.db
	//    [Web]
	//    port = 5555
	//    root = .speaks
	//    upload = data/store
	//    template = templates
	//    [LDAP]
	//    use = true
	//    BASEDN
	//    TEST
	//    [Session]
	//    secret -> generate
	//    name -> User???

	// args[1] init file...???

}

func main() {
	flag.Parse()

	args := flag.Args()
	leng := len(args)
	if leng != 1 {
		log.Println("Error: speaks arguments.")
		os.Exit(1)
	}

	sub := args[0]
	var err error

	switch sub {
	case "init":
		err = Init()
	case "start":
		err = Listen()
	case "help":
		err = Help()
	case "version":
		err = Version()
	default:
		log.Println("Error: speaks sub command(init | start | version | help)")
		os.Exit(1)
	}

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func Help() error {
	return fmt.Error("Not implements")
}

func Version() error {
	return fmt.Error("Not implements")
}

func Init() error {
	return fmt.Error("Not implements")
}

func Listen() error {

	iniFile := ".speaks/speaks.ini"
	err := Load(iniFile)
	if err != nil {
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
