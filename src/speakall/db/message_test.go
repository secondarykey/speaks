package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessage(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("message table test", t, func() {

		Convey("insert message table", func() {

			InsertMessage(10, "Category", "Content", "Created")
			rows, err := inst.Query("select user_id,category,content,created from Message")

			So(err, ShouldBeNil)
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

			msgs, err := SelectMessage("Category", "")

			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)
		})

		Convey("select message table", func() {
			_, err = InsertMessage(10, "Category", "Content", "Created")
			_, err = InsertMessage(11, "Category", "Content", "Created")
			_, err = InsertMessage(12, "TEST", "Content", "Created")

			msgs, err := SelectMessage("Category", "")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 2)

			msgs, err = SelectMessage("TEST", "")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)

			msgs, err = SelectMessage("", "")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 0)
		})
	})
}
