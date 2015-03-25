package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		teardown()
	}
	os.Exit(ret)
}

func setup() {
}

func teardown() {
}

func TestWeb(t *testing.T) {

	Convey("Listen", t, func() {
		Convey("Start WebServer", t, func() {
		})
	})

	Convey("setTemplates", t, func() {
		Convey("", t, func() {
		})
	})

	Convey("setJson", t, func() {
		Convey("", t, func() {
		})
	})
}
