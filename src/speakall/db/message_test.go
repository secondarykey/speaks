package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessage(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("database listen", t, func() {

		err = Listen("../../../data/db/test.db")
		So(err, ShouldBeNil)

		Convey("create message table", func() {
			err = createMessageTable()
			So(err, ShouldBeNil)
		})

		Convey("insert message table", func() {
			begin()
			InsertMessage(10, "Category", "Content")
			commit()

			rows, _ := inst.Query("select id,user_id,category,content,created from Message")
			defer rows.Close()
			for rows.Next() {
				var id int
				var user_id int
				var category string
				var content string
				var created string
				rows.Scan(&id, &user_id, &category, &content, &created)
				So(id, ShouldNotEqual, 0)
				So(user_id, ShouldEqual, 10)
				So(category, ShouldEqual, "Category")
				So(content, ShouldEqual, "Content")
				So(created, ShouldNotBeNil)
			}

			msgs, err := selectMessage("Category")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)
		})

		Convey("select message table", func() {
			begin()
			_, err = InsertMessage(10, "Category", "Content")
			_, err = InsertMessage(11, "Category", "Content")
			_, err = InsertMessage(12, "TEST", "Content")
			commit()

			msgs, err := selectMessage("Category")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 2)

			msgs, err = selectMessage("TEST")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)

			msgs, err = selectMessage("")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 0)
		})
	})
}
