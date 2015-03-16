package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type setting struct {
	Database database
	Web      web
	Session  session
}

type database struct {
	Path    string
	Version string
}

type web struct {
	Port   string
	Root   string
	Upload string
}

type session struct {
	Secret string
	Name   string
}

var Config setting

func init() {
	_, err := toml.DecodeFile("SpeakAll.ini", &Config)
	if err != nil {
		log.Println(err)
	}
}
