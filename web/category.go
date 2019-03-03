package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/secondarykey/speaks/db"
)

//category
func categoryHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	// /manage/category/
	if r.Method == "GET" {

		key, err := db.GenerateCategoryKey()
		if err != nil {
			return "", err
		}
		data["CategoryKey"] = key

		//TODO カテゴリリストを作成

		return "manage/category.tmpl", nil

	} else {

		name := r.FormValue("name")
		desc := r.FormValue("description")
		key := r.FormValue("key")

		u := data["User"].(*db.User)
		c := db.Category{
			Key:         key,
			Name:        name,
			Project:     u.CurrentProject.Key,
			Description: desc,
		}

		err := db.InsertCategory(c)
		if err != nil {
			return "", err
		}

		return "/manage/category/", NewRedirect("/manage/category/")
	}
}

func categoryListHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "GET" {
		return "", fmt.Errorf("NotSupport Method")
	}

	u := data["User"].(*db.User)
	p := u.CurrentProject

	cats, err := db.SelectProjectCategories(p.Key)
	if err != nil {
		return "", err
	}

	data["CategoryList"] = cats
	return "", nil
}

func categoryViewHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {
	if r.Method == "GET" {
		return "", fmt.Errorf("NotSupport Method")
	}

	url := r.URL.Path
	pathS := strings.Split(url, "/")

	cat, err := db.SelectCategory(pathS[3])
	if err != nil {
		return "", err
	}

	data["Category"] = cat
	return "", nil
}

func categoryDeleteHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "GET" {
		return "", fmt.Errorf("GET Method Error")
	}
	url := r.URL.Path
	pathS := strings.Split(url, "/")

	//TODO tx
	err := db.DeleteAllMessage(pathS[3])
	if err != nil {
		return "", err
	}
	err = db.DeleteCategory(pathS[3])
	if err != nil {
		return "", err
	}
	return "/maange/category/", NewRedirect("/manage/category/")
}
