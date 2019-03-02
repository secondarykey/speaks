package web

import (
	"io"
	"net/http"
	"os"
	"strings"

	. "github.com/secondarykey/speaks/config"

	uuid "github.com/satori/go.uuid"
)

type result struct {
	FileName string
}

func basePath() string {
	return Config.Base.Root + "/" + Config.Web.Upload
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	path := basePath()
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

	path := basePath()
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
