package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
	"github.com/secondarykey/speaks/http"
	"github.com/secondarykey/speaks/ws"
)

const Ver = "0.7.0"

var currentDir string
var logLv string
var logger *log.Logger

func init() {
	// -lv log : Level debug,fatal,error,
	flag.StringVar(&currentDir, "d", ".speaks", "CurrentDirectory")
	flag.StringVar(&logLv, "v", "WARN", "Log Level")

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func Usage() {

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
		err = Init(reader, currentDir)
	case "start":
		err = Start(currentDir)
	case "help":
		err = Help()
	case "version":
		err = Version()
	case "release":
		err = Release()
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

func Release() error {

	// delete ini file
	// delete db file

	// cd .speaks
	// call go-bindata

	//
	// go-bindata.exe -pkg=config -o=../config/binary.go .

	return fmt.Errorf("Not yet implemented.")
}

func Init(reader io.Reader, dir string) error {

	log.Println("Install Speaks[" + dir + "]")
	err := Ask(reader, dir)
	if err != nil {
		return err
	}
	return db.Init()
}

func Start(d string) error {

	log.Println("######## Initialize")
	err := Load(d)
	if err != nil {
		log.Println("Load Error:" + err.Error())
		return err
	}

	log.Println("######## start DBServer")
	err = db.Listen()
	if err != nil {
		log.Println("Database Listen:" + err.Error())
		return err
	}

	log.Println("######## start WSServer")
	err = ws.Listen("/ws/")
	if err != nil {
		log.Println("WebSocket Listen:" + err.Error())
		return err
	}

	log.Println("######## start HTTPServer")
	port := Config.Web.Port
	dir := Config.Base.Root

	err = http.Listen(dir, port)

	if err == nil {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for _ = range c {
			close(c)
		}
	}

	return err
}
