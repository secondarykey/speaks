package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
)

func switchHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	user := data["User"].(*db.User)

	//[/project/switch/{Key}]
	path := r.URL.Path
	pathS := strings.Split(path, "/")
	if len(pathS) != 4 {
		return "", fmt.Errorf("URL Error[%s].", path)
	}

	key := pathS[3]
	for _, elm := range user.Projects {
		if elm.Key == key {
			user.CurrentProject = elm
			saveLoginUser(r, w, user)
			return "/", NewRedirect("/")
		}
	}

	return "", fmt.Errorf("NotFound Project[%s].", key)
}

func userHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "POST" {
		//新規登録
		name := r.FormValue("dispName")
		email := r.FormValue("email")
		if name != "" && email != "" {
			u := &db.User{}
			u.Name = name
			u.Email = email

			err := db.CreateUser(u)
			if err != nil {
				return "", err
			}
		}
	}

	ulist, err := db.SelectAllUser()
	if err != nil {
		return "", err
	}

	for _, elm := range ulist {
		if strings.Contains(elm.Password, "-") {
			elm.Password = "/user/register/" + elm.Password
		} else {
			elm.Password = ""
		}
	}

	data["UserList"] = ulist

	return "admin/user.tmpl", nil
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//URLからユーザを特定
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	userKey := pathS[3]

	u, err := db.SelectPassword(userKey)
	if err != nil {
		return "", err
	}

	if r.Method == "GET" {
		data["User"] = u
		data["EditUser"] = u
		return "admin/user.tmpl", nil
	}

	//Form値を取得
	name := r.FormValue("dispName")
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	veri := r.FormValue("verifiedPassword")

	if email == "" {
		return "", fmt.Errorf("Empty Email")
	}

	if pwd == "" || pwd != veri {
		return "", fmt.Errorf("Bad Request")
	}

	u.Name = name
	u.Email = email
	u.Password = db.CreateMD5(pwd)

	err = db.UpdateUser(u)
	if err != nil {
		return "", err
	}

	return "/", NewRedirect("/")
}

func meHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	user := data["User"].(*db.User)
	data["EditUser"] = user

	if r.Method == "POST" {

		name := r.FormValue("dispName")
		email := r.FormValue("email")
		user.Name = name
		user.Email = email

		pwd := r.FormValue("password")
		veri := r.FormValue("verifiedPassword")
		if pwd != "" {
			if pwd == veri {
				user.Password = db.CreateMD5(pwd)
			}
		}

		//update
		err := db.UpdateUser(user)
		if err != nil {
			return "", err
		}

		err = saveLoginUser(r, w, user)
		if err != nil {
			return "", err
		}
	}

	return "admin/user.tmpl", nil
}

func iconRegisterHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	user := data["User"].(*db.User)
	//TODO static
	path := Config.Base.Root + "/static/images/icon/" + fmt.Sprint(user.Id)

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

	return "/", NewRedirect("/")
}
