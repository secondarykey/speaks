package config

import (
	"github.com/BurntSushi/toml"
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
	Port     string
	Root     string
	Upload   string
	Template string
}

type session struct {
	Secret string
	Name   string
}

var Config setting

func Load(file string) error {
	_, err := toml.DecodeFile(file, &Config)

	//Create path?

	return err
}
