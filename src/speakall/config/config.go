package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"strings"
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
	if err != nil {
		return err
	}

	paths := strings.Split(Config.Database.Path, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		//Create path?
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	return os.MkdirAll(Config.Web.Upload, 0777)
}
