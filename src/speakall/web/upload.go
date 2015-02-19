package web

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	. "speakall/config"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	path := Config.Web.Upload
	path += "/" + uuid.NewV4().String()

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

	http.Redirect(w, r, "/", http.StatusFound)
}
