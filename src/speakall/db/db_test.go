package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDB(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("database listen", t, func() {

		err = Listen("../../../data/db/SpeakAll-%s.db", "test")
		So(err, ShouldBeNil)

		Convey("Exec", func() {
			result, err := Exec("select * user")
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)

			result, err = Exec("select * from user")
			So(err, ShouldBeNil)
		})

		Convey("createInitTable", func() {
			// user exist
			// role exist
			// user_role exist
		})

		Convey("insertInitTable", func() {
			//user
			//role
			//user_role
		})
	})

}
