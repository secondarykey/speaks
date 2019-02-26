package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type setting struct {
	Base     base
	Database database
	Web      web
	LDAP     ldap
	Session  session
}

type base struct {
	Root string
}

type database struct {
	Path    string
	Version string
}

type web struct {
	Port     string
	Upload   string
	Template string
}

type ldap struct {
	Use      bool
	Server   string
	Protocol string
	Port     string
	BindDN   string
	BindPW   string
	BaseDn   string
	Filter   string
}

type session struct {
	Secret string
	Name   string
}

var Config *setting

func Load(file string) error {
	_, err := toml.DecodeFile(file, &Config)
	if err != nil {
		log.Println(err)
		return err
	}

	paths := strings.Split(Config.Database.Path, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		//Create path?
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return os.MkdirAll(Config.Web.Upload, 0777)
}

func Ask(reader io.Reader) error {

	conf := setting{}
	stdin := bufio.NewScanner(reader)
	var err error

	conf.Base.Root, err = ask(stdin, "Root Directory", ".speaks",
		func(in string) (string, error) {
			err := os.MkdirAll(in, 0777)
			if err != nil {
				fmt.Println(err)
				return "", nil
			}
			return in, nil
		})

	if err != nil {
		log.Println(err)
		return err
	}

	conf.Web.Port, err = ask(stdin, "HTTP Port Number", "5555",
		func(in string) (string, error) {
			_, err := strconv.ParseInt(in, 10, 64)
			if err != nil {
				fmt.Println(err)
				return "", nil
			}
			return in, nil
		})
	if err != nil {
		log.Println(err)
		return err
	}

	conf.Database.Version = "0.2"
	conf.Database.Path = "speaks-%s.db"

	conf.Web.Upload = "data/store"
	conf.Web.Template = ".speaks/templates"

	conf.Session.Secret = "UUID-AAA"
	conf.Session.Name = "User"

	//Secret
	//User
	Config = &conf
	return Config.Generate(".speaks/speaks.ini")
}

func ask(in *bufio.Scanner, msg string, def string, fn func(string) (string, error)) (string, error) {
	var err error
	for {
		fmt.Printf(msg+"[%s]:", def)
		in.Scan()

		input := in.Text()
		if input == "" {
			input = def
		}

		var datum string
		datum, err = fn(input)
		if err != nil {
			break
		}

		if datum == input {
			return input, nil
		}
	}

	log.Println(err)
	return "", fmt.Errorf("Input Error[%v]", err)
}

func (c *setting) Generate(f string) error {

	fs, err := os.Create(f)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fs.Close()

	fs.WriteString("[Base]\n")
	fs.WriteString(fmt.Sprintf("root=%s\n", c.Base.Root))
	fs.WriteString("\n")

	fs.WriteString("[Database]\n")
	fs.WriteString(fmt.Sprintf("version=%s\n", c.Database.Version))
	fs.WriteString(fmt.Sprintf("path=%s\n", c.Database.Path))
	fs.WriteString("\n")

	fs.WriteString("[Web]\n")
	fs.WriteString(fmt.Sprintf("port=%s\n", c.Web.Port))
	fs.WriteString(fmt.Sprintf("upload=%s\n", c.Web.Upload))
	fs.WriteString(fmt.Sprintf("template=%s\n", c.Web.Template))
	fs.WriteString("\n")

	fs.WriteString("[LDAP]\n")
	fs.WriteString(fmt.Sprintf("use=%t\n", c.LDAP.Use))

	fs.WriteString("\n")

	fs.WriteString("[Session]\n")
	fs.WriteString(fmt.Sprintf("Secret=%s\n", c.Session.Secret))
	fs.WriteString(fmt.Sprintf("name=%s\n", c.Session.Name))
	fs.WriteString("\n")

	return nil
}
