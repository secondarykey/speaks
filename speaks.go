package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
	"github.com/secondarykey/speaks/web"
	"github.com/secondarykey/speaks/ws"
)

const Ver = "0.5.0"

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}

func Usage() {

	// -lv log : Level debug,fatal,error,
	// args[0] sub command : init start version help

}

func main() {

	flag.Parse()
	args := flag.Args()

	err := run(os.Stdin, args)
	if err != nil {
		log.Printf("Error:%v\n", err)
		os.Exit(1)
	}

	log.Println("Success:speaks.")
	os.Exit(0)
}

func run(reader io.Reader, args []string) error {

	leng := len(args)
	if leng != 1 {
		return fmt.Errorf("speaks Agument required sub command.")
	}

	sub := args[0]
	var err error

	switch sub {
	case "init":
		err = Init(reader)
	case "start":
		err = Start(".speaks/speaks.ini")
	case "help":
		err = Help()
	case "version":
		err = Version()
	default:
		return fmt.Errorf("Error: speaks sub command(init | start | version | help)")
	}

	if err != nil {
		return err
	}

	return nil
}

func Help() error {
	Usage()
	return nil
}

func Version() error {
	log.Printf("Speaks Version %s\n", Ver)
	return nil
}

func Init(reader io.Reader) error {
	log.Println("Install Speaks[.speak]")
	err := Ask(reader)
	if err != nil {
		return err
	}
	return nil
}

func Start(f string) error {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("シグナル来た", sig)
			close(c)
			os.Exit(0)
		}
	}()

	log.Println("######## Initialize")
	err := Load(f)
	if err != nil {
		return err
	}

	log.Println("######## start DBServer")
	path := Config.Database.Path
	ver := Config.Database.Version
	err = db.Listen(path, ver)
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
	dir := Config.Base.Root

	return web.Listen(dir, port)
}
