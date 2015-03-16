package speakall

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

func TestSpeakAll(t *testing.T) {

	Convey("speakall start", t, func() {

		Convey("database", t, func() {
		})

		Convey("websocket", t, func() {
		})

		Convey("web", t, func() {
		})
	})
}
