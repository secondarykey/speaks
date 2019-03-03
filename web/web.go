package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
)

func init() {

	//All
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	//Speaker
	http.HandleFunc("/", chatHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/store/", storeHandler)

	http.HandleFunc("/message/delete/", messageDeleteHandler)

	//Editor
	http.HandleFunc("/memo", memoListHandler)
	http.HandleFunc("/memo/edit/", memoEditHandler)
	http.HandleFunc("/memo/view/", memoViewHandler)
}

type handler func(http.ResponseWriter, *http.Request, map[string]interface{}) (string, error)
type urlPattern map[string]handler

type router struct {
	pattern urlPattern
}

type htmlRouter router
type jsonRouter router

func NewURLPattern() urlPattern {
	inst := make(map[string]handler)
	return urlPattern(inst)
}

func (p urlPattern) getHandle(path string) (handler, bool) {
	hand, ok := p[path]
	if ok {
		return hand, true
	}
	for key, elm := range p {
		if strings.Index(path, key) == 0 {
			return elm, true
		}
	}
	return hand, false
}

func NewHTMLRouter() *htmlRouter {

	router := htmlRouter{}
	router.pattern = NewURLPattern()

	//All
	router.pattern["/me"] = meHandler
	router.pattern["/me/upload"] = iconRegisterHandler
	router.pattern["/project/switch/"] = switchHandler

	//Manage
	router.pattern["/manage/category/"] = categoryHandler
	router.pattern["/manage/category/delete/"] = categoryDeleteHandler
	router.pattern["/manage/project/member"] = projectHandler

	//Admin
	router.pattern["/admin/project/"] = projectHandler
	router.pattern["/admin/user"] = userHandler
	router.pattern["/admin/user/register"] = userRegisterHandler

	return &router
}

func (route htmlRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//login,logout は使わない

	//Login
	u, err := getLoginUser(r)
	if err != nil {
		//Redirect
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	obj := make(map[string]interface{})
	obj["User"] = u

	path := r.URL.Path
	obj["URL"] = path
	obj["Type"] = NewType(path)

	tmpl := ""
	log.Println("Path:" + path)

	handle, ok := route.pattern.getHandle(path)
	if ok {
		tmpl, err = handle(w, r, obj)
	} else {
		err = fmt.Errorf("NotFound Handler[%s]", path)
	}

	if err != nil {
		if rd, ok := err.(*Redirect); ok {
			http.Redirect(w, r, string(rd.Path()), http.StatusFound)
			return
		}

		tmpl = "error.tmpl"
		obj["Error"] = err
		log.Println(err)
	}

	//Template Write
	err = setTemplates(w, obj, tmpl)
	if err != nil {
		log.Println(err)
	}
}

func NewJSONRouter() *jsonRouter {
	router := jsonRouter{}
	router.pattern = NewURLPattern()

	router.pattern["/api/message/"] = messageHandler

	router.pattern["/api/category/list"] = categoryListHandler
	router.pattern["/api/category/view/"] = categoryViewHandler
	return &router
}

func (route jsonRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var err error
	var u *db.User

	path := r.URL.Path
	log.Println("Path:" + path)

	data := make(map[string]interface{})
	//Login
	u, err = getLoginUser(r)

	if err == nil {
		data["User"] = u
		path := r.URL.Path
		handle, ok := route.pattern.getHandle(path)
		if ok {
			_, err = handle(w, r, data)
		} else {
			err = fmt.Errorf("NotFound Handler[%s]", path)
		}
	}

	if err != nil {
		log.Println(err)
		data["Error"] = err.Error()
	}

	err = setJson(data, w)
	if err != nil {
		log.Println(err)
	}
}

type templateType string

func NewType(path string) templateType {
	return templateType(path)
}

func (t templateType) IsAdmin() bool {
	return t.prefixPattern("/admin")
}

func (t templateType) IsManage() bool {
	return t.prefixPattern("/manage")
}

func (t templateType) IsMe() bool {
	return t.prefixPattern("/me")
}

func (t templateType) IsChat() bool {
	return t.prefixPattern("/")
}

func (t templateType) prefixPattern(prefix string) bool {
	if strings.Index(string(t), prefix) == 0 {
		return true
	}
	return false
}

type Redirect string

func NewRedirect(path string) *Redirect {
	r := Redirect(path)
	return &r
}

func (r *Redirect) Error() string {
	return "Redirect type is not Error."
}

func (r *Redirect) Path() string {
	return string(*r)
}

func Listen(static, port string) error {

	startSession()
	//TODO Static
	http.Handle("/static/", http.FileServer(http.Dir(static)))

	router := NewHTMLRouter()
	http.Handle("/category/", router)
	http.Handle("/project/", router)
	http.Handle("/admin/", router)
	http.Handle("/manage/", router)
	http.Handle("/me", router)

	api := NewJSONRouter()
	http.Handle("/api/", api)

	return http.ListenAndServe(":"+port, nil)
}

func setTemplates(w http.ResponseWriter, param interface{}, templateFile string) error {

	templateDir := Config.Base.Root + "/" + Config.Web.Template
	tmpls := make([]string, 0)
	tmpls = append(tmpls, templateDir+"/layout.tmpl")
	tmpls = append(tmpls, templateDir+"/menu.tmpl")
	tmpls = append(tmpls, templateDir+"/"+templateFile)

	tmpl := template.Must(template.ParseFiles(tmpls...))
	if err := tmpl.Execute(w, param); err != nil {
		return err
	}
	return nil
}

func setJson(s interface{}, w http.ResponseWriter) error {
	res, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	return err
}
