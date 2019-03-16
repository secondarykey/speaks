package http

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

func storeHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {
	path := basePath()
	urlSlice := strings.Split(r.URL.Path, "/")

	file := path + "/" + urlSlice[2]
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.ServeFile(w, r, file)
	}
	return "", NewNoWrite("Write Binary")
}

func uploadHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	path := basePath()
	fileId := uuid.NewV4().String()
	path += "/" + fileId

	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		return "", err
	}
	defer file.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	ret := &result{
		FileName: "store/" + fileId,
	}

	data["Result"] = ret
	return "", nil
}
