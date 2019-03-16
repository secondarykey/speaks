package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/secondarykey/speaks/db"
)

func memberHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//ユーザを作成
	users, err := db.SelectAllUser()
	if err != nil {
		return "", err
	}

	u := data["User"].(*db.User)
	key := u.CurrentProject.Key

	member, err := db.SelectProjectMember(key)
	if err != nil {
		return "", err
	}

	for _, mem := range member {
		var target *db.User
		for _, elm := range users {
			if mem.UserId == elm.Id {
				target = elm
				break
			}
		}

		if target == nil {
			log.Println("Error no target[%d]", mem.UserId)
			continue
		}
		target.CurrentProject = u.CurrentProject

		pr := target.ProjectRoles
		if pr == nil {
			pr = make(map[string]db.RoleMap)
		}

		rm := pr[key]
		if rm == nil {
			rm = db.NewMemberRole()
		}

		rm[mem.Role] = true

		pr[key] = rm
		target.ProjectRoles = pr
	}

	//全ユーザのロールを設定

	data["UserList"] = users

	return "manage/member.tmpl", nil
}

func memberUpdateHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//カンマ区切り(ロール名-ユーザID)
	roles := r.FormValue("roleMember")
	slice := strings.Split(roles, ",")

	u := data["User"].(*db.User)
	key := u.CurrentProject.Key

	//TODO tx
	//全メンバの削除
	err := db.DeleteProjectMembers(key)
	if err != nil {
		return "", err
	}

	var members []db.Member
	//全メンバの追加
	for _, v := range slice {
		roleId := strings.Split(v, "-")
		m := db.Member{}
		m.Role = roleId[0]
		m.UserId, _ = strconv.Atoi(roleId[1])
		m.Project = key
		members = append(members, m)
	}

	err = db.InsertMembers(members)
	if err != nil {
		return "", err
	}

	return "/manage/project/member", NewRedirect("/manage/project/member")
}
