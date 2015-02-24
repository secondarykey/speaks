package speakall

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSpeakAll(t *testing.T) {

	Convey("speakall start", t, func() {

		go Start()

		Convey("database", t, func() {
		})
		Convey("websocket", t, func() {
		})
		Convey("web", t, func() {
		})
	})
}
