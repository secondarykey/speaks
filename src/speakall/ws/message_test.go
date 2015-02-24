package ws

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWSMessage(t *testing.T) {
	Convey("createOpenMessage", t, func() {
		ms := createOpenMessage("clientId")
		Convey("Type", t, func() {
			So(ms.Type, ShouldEqual, "Open")
			So(ms.ClientId, ShouldEqual, "clientId")
			So(ms.Content, ShouldEqual, "clientId")
		})
	})
}
