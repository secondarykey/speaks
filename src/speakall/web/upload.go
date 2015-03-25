package web

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	. "speakall/config"
	"strings"
)

type result struct {
	FileName string
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	path := Config.Web.Upload
	urlSlice := strings.Split(r.URL.Path, "/")

	file := path + "/" + urlSlice[2]
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.ServeFile(w, r, file)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	path := Config.Web.Upload
	fileId := uuid.NewV4().String()
	path += "/" + fileId

	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	out, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret := &result{
		FileName: "store/" + fileId,
	}

	setJson(ret, w)
}
