package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSession(t *testing.T) {

	Convey("session test", t, func() {

		So(store, ShouldBeNil)

		Convey("startSession", t, func() {
			startSession("Secret")
			So(store, ShouldNotBeNil)
		})

		Convey("getSession", t, func() {
		})
	})
}
