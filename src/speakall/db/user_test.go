package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUser(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("create user table", t, func() {

		err = Listen("../../../data/db/test.db")
		So(err, ShouldBeNil)

		Convey("insert user table", func() {
		})
		Convey("createMD5", func() {
		})
		Convey("select user", func() {
		})
	})

}
