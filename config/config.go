package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	statik "github.com/rakyll/statik/fs"
	uuid "github.com/satori/go.uuid"
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
	Public   string
}

type ldap struct {
	Use          bool
	Host         string
	BaseDN       string
	BindDN       string
	BindPassword string
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
	return nil
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
	conf.Base.Root = root

	//[Database]
	conf.Database.Version = "1.0.0"
	conf.Database.Path = "speaks-%s.db"

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

	conf.Web.Port = port
	conf.Web.Upload = "data"
	conf.Web.Template = "templates"
	conf.Web.Public = "public"

	//[LDAP]
	conf.LDAP.Use = false
	ldap, err := ask(stdin, "Use LDAP(ActiveDirectory)?", "Y/n",
		func(in string) (string, error) {
			return in, nil
		})
	if err != nil {
		log.Println(err)
		return err
	}
	if ldap == "Y" {
		err = askLDAP(&conf, stdin)
		if err != nil {
			return err
		}
	}

	secret := uuid.NewV4().String()
	conf.Session.Secret = secret
	conf.Session.Name = "User050"

	//[Session]
	//Secret
	//User
	Config = &conf
	return Config.Generate(root, DefaultInitFileName)
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

func askLDAP(c *setting, stdin *bufio.Scanner) error {

	c.LDAP.Use = true

	c.LDAP.Host = "localhost"
	c.LDAP.BaseDN = "dc=sample,dc=com"
	c.LDAP.BindDN = "user@sample.com"
	c.LDAP.BindPassword = "***"

	host, err := ask(stdin, "LDAP Host", c.LDAP.Host,
		func(in string) (string, error) { return in, nil })
	if err != nil {
		return err
	}
	c.LDAP.Host = host

	base, err := ask(stdin, "LDAP BaseDN", c.LDAP.BaseDN,
		func(in string) (string, error) { return in, nil })
	if err != nil {
		return err
	}
	c.LDAP.BaseDN = base

	dn, err := ask(stdin, "LDAP Bind User", c.LDAP.BindDN,
		func(in string) (string, error) { return in, nil })
	if err != nil {
		return err
	}
	c.LDAP.BindDN = dn

	pwd, err := ask(stdin, "LDAP Bind User Password", c.LDAP.BindPassword,
		func(in string) (string, error) { return in, nil })
	if err != nil {
		return err
	}
	c.LDAP.BindPassword = pwd

	return nil
}

func (c *setting) Generate(d, f string) error {

	path := d + "/" + f
	fs, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fs.Close()

	fs.WriteString("[Base]\n")
	fs.WriteString(fmt.Sprintf(`root="%s"`+"\n", c.Base.Root))
	fs.WriteString("\n")

	fs.WriteString("[Database]\n")
	fs.WriteString(fmt.Sprintf(`version="%s"`+"\n", c.Database.Version))
	fs.WriteString(fmt.Sprintf(`path="%s"`+"\n", c.Database.Path))
	fs.WriteString("\n")

	fs.WriteString("[Web]\n")
	fs.WriteString(fmt.Sprintf(`port="%s"`+"\n", c.Web.Port))
	fs.WriteString(fmt.Sprintf(`upload="%s"`+"\n", c.Web.Upload))
	fs.WriteString(fmt.Sprintf(`template="%s"`+"\n", c.Web.Template))
	fs.WriteString(fmt.Sprintf(`public="%s"`+"\n", c.Web.Public))
	fs.WriteString("\n")

	fs.WriteString("[LDAP]\n")
	fs.WriteString(fmt.Sprintf("use=%t\n", c.LDAP.Use))
	fs.WriteString(fmt.Sprintf(`host="%s"`+"\n", c.LDAP.Host))
	fs.WriteString(fmt.Sprintf(`baseDN="%s"`+"\n", c.LDAP.BaseDN))
	fs.WriteString(fmt.Sprintf(`bindDN="%s"`+"\n", c.LDAP.BindDN))
	fs.WriteString(fmt.Sprintf(`bindPassword="%s"`+"\n", c.LDAP.BindPassword))

	fs.WriteString("\n")

	fs.WriteString("[Session]\n")
	fs.WriteString(fmt.Sprintf(`Secret="%s"`+"\n", c.Session.Secret))
	fs.WriteString(fmt.Sprintf(`name="%s"`+"\n", c.Session.Name))
	fs.WriteString("\n")

	log.Println("Generate speaks.ini")

	statikFS, err := statik.New()
	if err != nil {
		return err
	}

	statik.Walk(statikFS, "/", func(path string, info os.FileInfo, err error) error {

		name := d + path

		log.Println(name)

		if info.IsDir() {
			os.MkdirAll(name, 0777)
		} else {

			bf, err := os.Create(name)
			if err != nil {
				log.Println(err)
				return err
			}
			defer bf.Close()

			fp, err := statikFS.Open(path)
			if err != nil {
				log.Println(err)
				return err
			}
			defer fp.Close()

			_, err = io.Copy(bf, fp)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		return nil
	})

	dataDir := d + "/" + Config.Web.Upload
	os.MkdirAll(dataDir, 0777)

	dataDir = d + "/" + Config.Web.Public + "/images/icon"
	os.MkdirAll(dataDir, 0777)

	log.Println("Generate static file")

	return nil
}
