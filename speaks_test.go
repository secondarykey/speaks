package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {

	args := make([]string, 0)
	stdin := strings.NewReader("\n\n")

	err := run(stdin, args)
	if err == nil {
		t.Error("Argument error Nil")
	}

	args = append(args, "error")
	err = run(stdin, args)
	if err == nil {
		t.Error("Argument sub command error Nil")
	}

	args[0] = "init"
	err = run(stdin, args)
	if err != nil {
		t.Errorf("run() Init error not Nil[%v]", err)
	}

	args[0] = "version"
	err = run(stdin, args)
	if err != nil {
		t.Errorf("run() version error not Nil[%v]", err)
	}

	args[0] = "help"
	err = run(stdin, args)
	if err != nil {
		t.Errorf("run() help error not Nil[%v]", err)
	}

	//TODO start

}
