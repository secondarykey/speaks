package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
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
}

func teardown() {
}

func TestWeb(t *testing.T) {

	Convey("Listen", t, func() {
		Convey("database", t, func() {
		})
	})

	Convey("handler", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		res, err := http.Get(ts.URL)
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("meHandler", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(meHandler))
		defer ts.Close()
		res, err := http.Get(ts.URL)
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

}
