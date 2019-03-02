package config

import (
	"bufio"
	"bytes"
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

const DefaultInitFileName = "speaks.ini"

func Load(d string) error {

	file := d + "/" + DefaultInitFileName
	_, err := toml.DecodeFile(file, &Config)
	if err != nil {
		log.Println(err)
		return err
	}

	//All Create
	paths := strings.Split(Config.Database.Path, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return os.MkdirAll(Config.Web.Upload, 0777)
}

func Ask(reader io.Reader, root string) error {

	conf := setting{}

	stdin := bufio.NewScanner(reader)
	var err error

	//[Base]
	err = os.MkdirAll(root, 0777)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	conf.Base.Root = `"` + root + `"`

	//[Database]
	conf.Database.Version = `"1.0.0"`
	conf.Database.Path = `"speaks-%s.db"`

	//[Web]
	port, err := ask(stdin, "HTTP Port Number", "5555",
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

	conf.Web.Port = `"` + port + `"`
	conf.Web.Upload = `"data/store"`
	conf.Web.Template = `"templates"`

	//TODO LDAP Ask
	conf.LDAP.Use = false

	//TODO UUID
	conf.Session.Secret = `"UUID-AAA"`
	conf.Session.Name = `"User"`

	//Secret
	//User
	Config = &conf
	return Config.Generate(root + "/" + DefaultInitFileName)
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

	log.Println("Generate speaks.ini")
	names := AssetNames()

	//TODO method
	os.MkdirAll(".speaks/static/js", 0777)
	os.MkdirAll(".speaks/static/images/icon", 0777)
	os.MkdirAll(".speaks/static/css", 0777)
	os.MkdirAll(".speaks/data/store", 0777)
	os.MkdirAll(".speaks/templates/memo", 0777)

	for _, name := range names {

		bin, err := Asset(name)
		if err != nil {
			log.Println(err)
			return err
		}

		reader := bytes.NewReader(bin)
		bf, err := os.Create(name)
		if err != nil {
			log.Println(err)
			return err
		}
		defer bf.Close()

		_, err = io.Copy(bf, reader)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("Generate " + name)
	}

	return nil
}
