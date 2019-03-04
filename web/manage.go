package web

import (
	"log"
	"net/http"

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
			log.Println("Error no target")
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
