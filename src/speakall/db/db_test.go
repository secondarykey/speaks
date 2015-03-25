package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setUp()
	ret := m.Run()
	if ret == 0 {
		tearDown()
	}
	os.Exit(ret)
}

func setUp() {
	err := Listen("../../../data/db/SpeakAll-%s.db", "test")
	if err != nil {
		panic(err)
	}
}

func tearDown() {
}

func ListenTestTable() {
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

		Convey("createInitTables", func() {
			// user exist
			// role exist
			// user_role exist
		})

		Convey("deleteTables", func() {
			//user
			//role
			//user_role
		})

		Convey("insertInitTable", func() {
			//user
			//role
			//user_role
		})
	})

}
