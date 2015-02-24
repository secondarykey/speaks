package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {

	Convey("Listen", t, func() {
		Convey("database", t, func() {
		})
	})

	Convey("handler", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		res, err := http.Get("/")
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("me", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(meHandler))
		defer ts.Close()
		res, err := http.Get("/me")
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

}
