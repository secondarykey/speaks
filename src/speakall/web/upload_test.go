package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpload(t *testing.T) {

	Convey("uploadHandler", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(uploadHandler))
		defer ts.Close()
		res, err := http.Get(ts.URL)
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("storeHandler", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(storeHandler))
		defer ts.Close()
		res, err := http.Get(ts.URL)
		So(err, ShouldBeNil)

		Convey("status code", t, func() {
			So(res.StatusCode, ShouldEqual, 404)
		})
	})
}
