package db

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
	err := Listen("../../../data/db/SpeakAll-%s.db", "test")
	So(err, ShouldBeNil)
}

func teardown() {
}

func TestDB(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("database test", t, func() {

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
		Convey("begin", func() {
		})
		Convey("rollback", func() {
		})
		Convey("commit", func() {
		})
	})

}
