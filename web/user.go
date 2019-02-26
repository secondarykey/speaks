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

func userHandler(w http.ResponseWriter, r *http.Request) {

	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	ulist, err := db.SelectAllUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tc := make(map[string]interface{})
	tc["URL"] = "/user"
	for _, elm := range ulist {
		if strings.Contains(elm.Password, "-") {
			elm.Password = "/user/register/" + elm.Password
		} else {
			elm.Password = ""
		}
	}

	tc["User"] = user
	tc["UserList"] = ulist

	setTemplates(w, tc, "menu.tmpl", "user.tmpl")
	return
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {

	//URLからユーザを特定
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	userKey := pathS[3]

	u, err := db.SelectPassword(userKey)
	if err != nil {
		http.Error(w, "Bad URL", http.StatusNotFound)
		return
	}

	//パスワードからユーザを設定
	tc := make(map[string]interface{})
	tc["User"] = u
	tc["URL"] = url

	if r.Method == "GET" {
		setTemplates(w, tc, "menu.tmpl", "user.tmpl")
		return
	}

	//Form値を取得
	name := r.FormValue("dispName")
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	veri := r.FormValue("verifiedPassword")

	if email == "" {
		http.Error(w, "Empty Email", http.StatusBadRequest)
		return
	}

	if pwd == "" || pwd != veri {
		http.Error(w, "Bad Password", http.StatusBadRequest)
		return
	}

	u.Name = name
	u.Email = email
	u.Password = db.CreateMD5(pwd)

	err = db.UpdateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func meHandler(w http.ResponseWriter, r *http.Request) {

	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["EditUser"] = user
	tc["URL"] = "/me"

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = saveLoginUser(r, w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	setTemplates(w, tc, "menu.tmpl", "user.tmpl")
}

func iconRegisterHandler(w http.ResponseWriter, r *http.Request) {

	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	path := Config.Base.Root + "/static/images/icon/" + fmt.Sprint(user.Id)

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
