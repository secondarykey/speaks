package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessage(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("message table test", t, func() {

		Convey("create message table", func() {
			err = createMessageTable()
			So(err, ShouldBeNil)
		})

		Convey("insert message table", func() {
			begin()
			InsertMessage(10, "Category", "Content", "Created")
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

			msgs, err := SelectMessage("Category", "")
			So(err, ShouldBeNil)
			So(len(msgs), ShouldEqual, 1)
		})

		Convey("select message table", func() {
			begin()
			_, err = InsertMessage(10, "Category", "Content", "Created")
			_, err = InsertMessage(11, "Category", "Content", "Created")
			_, err = InsertMessage(12, "TEST", "Content", "Created")
			commit()

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
