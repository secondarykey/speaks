package http

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

	leng := 0
	rtn := false
	for key, elm := range p {
		if strings.Index(path, key) == 0 {
			if len(key) > leng {
				hand = elm
				rtn = true
				leng = len(key)
				log.Println("Pattern:" + key)
			}
		}
	}
	return hand, rtn
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
	router.pattern["/manage/project/member"] = memberHandler
	router.pattern["/manage/project/member/update"] = memberUpdateHandler

	//Admin
	router.pattern["/admin/project/"] = projectHandler
	router.pattern["/admin/user"] = userHandler
	//router.pattern["/admin/user/register"] = userRegisterHandler

	//All
	router.pattern["/login"] = loginHandler
	router.pattern["/logout"] = logoutHandler

	//Speaker
	router.pattern["/"] = chatHandler
	router.pattern["/store/"] = storeHandler

	//Editor
	router.pattern["/memo"] = memoListHandler
	router.pattern["/memo/edit/"] = memoHandler
	router.pattern["/memo/"] = memoViewHandler

	router.pattern["/user/register/"] = userRegisterHandler
	return &router
}

func ignore(path string) bool {
	if strings.Index(path, "/user/register/") == 0 ||
		strings.Index(path, "/login") == 0 {
		return true
	}

	return false
}

func (route htmlRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	u, err := getLoginUser(r)

	if !ignore(path) {
		//Login
		if err != nil {
			//Redirect
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	} else if path == "/login" {
		u = new(db.User)
	}

	obj := make(map[string]interface{})
	obj["User"] = u

	obj["URL"] = path
	obj["Type"] = NewType(path)

	tmpl := ""

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

		if _, ok := err.(*NoWrite); ok {
			log.Println(err.Error())
			return
		}

		tmpl = "error.tmpl"
		obj["Error"] = err

		log.Printf("Path[%s]:%s\n", path, err)
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

	router.pattern["/api/message/delete/"] = messageDeleteHandler
	router.pattern["/api/message/"] = messageHandler

	router.pattern["/api/category/list"] = categoryListHandler
	router.pattern["/api/category/view/"] = categoryViewHandler

	router.pattern["/api/upload"] = uploadHandler

	router.pattern["/api/memo/delete/"] = memoDeleteHandler
	router.pattern["/api/memo/edit/"] = memoUpdateHandler
	return &router
}

func (route jsonRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var err error
	var u *db.User

	path := r.URL.Path

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
		log.Printf("Path[%s]:%s\n", path, err)
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

type NoWrite string
type Redirect string

func NewNoWrite(path string) *NoWrite {
	r := NoWrite(path)
	return &r
}
func (r *NoWrite) Error() string {
	return string(*r)
}

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

func Listen(root, port string) error {

	startSession()

	pub := http.Dir(root + "/" + Config.Web.Public)
	s := http.FileServer(pub)
	http.Handle("/js/", s)
	http.Handle("/css/", s)
	http.Handle("/images/", s)
	http.HandleFunc("/favicon.ico", faviconHandler)

	api := NewJSONRouter()
	http.Handle("/api/", api)

	router := NewHTMLRouter()
	http.Handle("/", router)

	return http.ListenAndServe(":"+port, nil)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s/%s/images/favicon.ico", Config.Base.Root, Config.Web.Public)
	http.ServeFile(w, r, path)
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
