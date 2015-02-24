package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	Convey("loginHandler start", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(loginHandler))
		defer ts.Close()
		res, err := http.Get("/login")
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("logoutHandler start", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(logoutHandler))
		defer ts.Close()
		res, err := http.Get("/logout")
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})
}
