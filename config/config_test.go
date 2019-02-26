package config

import (
	"os"
	"strings"
	"testing"
)

func TestAsk(t *testing.T) {
	s, err := Ask(strings.NewReader("\n\n"))
	if err != nil {
		t.Error("Not Error")
	}

	if s.Root != ".speaks" {
		t.Error("Not Error")
	}

	if s.Web.Port != "5555" {
		t.Error("Not Error")
	}
}

func TestGenerate(t *testing.T) {

	s := setting{}
	err := s.Generate("test.ini")
	if err != nil {
		t.Errorf("Generate() return not nil.[%v]", err)
	}

	_, err = os.Stat("test.ini")
	if err != nil {
		t.Errorf("test.ini not exist?[%v]", err)
	}

}

func TestLoad(t *testing.T) {

}
