package speakall

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	. "speakall/config"
	"testing"
)

func TestMain(m *testing.M) {
	err := Load("../../test.ini")
	if err != nil {
		panic(err)
	}

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
	})
}
