package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/secondarykey/speaks/db"
)

// /manage/category/
func categoryHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	u := data["User"].(*db.User)
	p := u.CurrentProject

	if r.Method == "GET" {

		key, err := db.GenerateCategoryKey()
		if err != nil {
			return "", err
		}
		data["CategoryKey"] = key

		cats, err := db.SelectProjectCategories(p.Key)
		if err != nil {
			return "", err
		}
		data["CategoryList"] = cats

		return "manage/category.tmpl", nil

	} else {

		name := r.FormValue("name")
		desc := r.FormValue("description")
		key := r.FormValue("key")

		c := db.Category{
			Key:         key,
			Name:        name,
			Project:     p.Key,
			Description: desc,
		}

		err := db.InsertCategory(c)
		if err != nil {
			return "", err
		}

		return "/manage/category/", NewRedirect("/manage/category/")
	}
}

// /api/category/list
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

// /api/category/view/{id}
func categoryViewHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "GET" {
		return "", fmt.Errorf("Method Error")
	}

	url := r.URL.Path
	pathS := strings.Split(url, "/")
	if len(pathS) != 5 {
		return "", fmt.Errorf("URL Error[%s]", url)
	}

	u := data["User"].(*db.User)
	key := pathS[4]

	cat, err := db.SelectCategory(key, u.CurrentProject.Key)
	if err != nil {
		return "", err
	}

	data["Category"] = cat
	return "", nil
}

// /manage/category/delete/{id}
func categoryDeleteHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	url := r.URL.Path
	pathS := strings.Split(url, "/")
	if len(pathS) != 5 {
		return "", fmt.Errorf("URL Error[%s]", url)
	}
	cat := pathS[4]

	u := data["User"].(*db.User)
	pro := u.CurrentProject

	err := db.DeleteCategory(pro.Key, cat)
	if err != nil {
		return "", err
	}
	return "/maange/category/", NewRedirect("/manage/category/")
}
