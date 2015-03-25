package main

import (
	"flag"
	"log"
	"speakall"
	"speakall/config"
)

func main() {

	var iniFile string
	//flag.StringVar(&iniFile, "ini", "SpeakAll.ini", "initialize file path.default 'SpeakAll.ini'")
	flag.Parse()
	args := flag.Args()
	leng := len(args)

	switch leng {
	case 0:
		iniFile = "SpeakAll.ini"
	case 1:
		iniFile = args[0]
	default:
		log.Println("Error: too many argument.")
		return
	}

	err := config.Load(iniFile)
	if err != nil {
		log.Println(err.Error())
		log.Println("Error: Loading initialize file.[" + iniFile + "]")
		return
	}

	err = speakall.Start()
	if err != nil {
		log.Println(err.Error())
		log.Println("Error: Start SpeakAll Server")
	}
}
