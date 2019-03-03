package web

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		setTemplates(w, nil, "login.tmpl")
		return
	}

	email := r.FormValue("email")
	pswd := r.FormValue("password")

	var err error
	var user *db.User

	if Config.LDAP.Use {
		user, err = loginLDAP(email, pswd)
	} else {
		user, err = loginDB(email, pswd)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = setUserRole(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = setProjectRole(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = setCurrentProject(user, "Speaks")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveLoginUser(r, w, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//ERROR PAGES
func loginLDAP(id, pass string) (*db.User, error) {
	return nil, fmt.Errorf("LDAP Login Not yet Implemented.")
}

func loginDB(id, pass string) (*db.User, error) {
	return db.SelectUser(id, pass)
}

func setUserRole(u *db.User) error {
	roles, err := db.SelectUserRole(u.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	u.Roles = db.NewRoleMap()
	for _, elm := range roles {
		u.Roles[elm.RoleKey] = true
	}
	return nil
}

func setProjectRole(u *db.User) error {

	projects, err := db.SelectProjects()
	if err != nil {
		return err
	}

	members, err := db.SelectMember(u.Id)
	if err != nil {
		return err
	}

	u.ProjectRoles = make(map[string]db.RoleMap)
	for _, elm := range members {
		key := elm.Project
		role, ok := u.ProjectRoles[key]
		if !ok {
			role = db.NewMemberRole()
			u.ProjectRoles[key] = role
		}
		role[elm.Role] = true
	}

	u.Projects = make([]db.Project, 0)
	for _, p := range projects {
		if u.See(p.Key) {
			u.Projects = append(u.Projects, p)
		}
	}
	return nil
}

func setCurrentProject(u *db.User, key string) error {

	if !u.See(key) {
		return fmt.Errorf("Auth Error")
	}

	for _, p := range u.Projects {
		if p.Key == key {
			u.CurrentProject = p
			return nil
		}
	}
	return fmt.Errorf("Not Found")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := saveLoginUser(r, w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
