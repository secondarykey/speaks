package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {

	Convey("config", t, func() {
		So(Config.Database, ShouldNotBeNil)
		So(Config.Database.Path, ShouldNotBeNil)
		So(Config.Database.Version, ShouldNotBeNil)
		So(Config.Web, ShouldNotBeNil)
		So(Config.Web.Port, ShouldNotBeNil)
		So(Config.Web.Root, ShouldNotBeNil)
		So(Config.Web.Upload, ShouldNotBeNil)
		So(Config.Session, ShouldNotBeNil)
		So(Config.Session.Secret, ShouldNotBeNil)
		So(Config.Session.Name, ShouldNotBeNil)
	})
}
