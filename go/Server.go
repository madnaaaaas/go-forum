package main

import (
	"net/http"
	//"net/http/cookiejar"
	"crypto/md5"
	"fmt"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetMd5(Value string) string {
	s := ""
	for _, v := range md5.Sum([]byte(Value)) {
		s += strconv.FormatUint(uint64(v), 16)
	}
	return s
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ind := GetIndexPage()
	ind.Own, ind.Log = TestLog(r)
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, ind)
}

func SubforumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/subforum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckSubforum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	sub := GetSubforumPage(Id)
	sub.Own, sub.Log = TestLog(r)
	t, _ := template.ParseFiles("html/subforum.html")
	t.Execute(w, sub)
}

func ThemeHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/theme/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckTheme(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	th := GetThemePage(Id)
	th.Own, th.Log = TestLog(r)
	t, _ := template.ParseFiles("html/theme.html")
	t.Execute(w, th)
}

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/forum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckForum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	f := GetForumPage(Id)
	f.Own, f.Log = TestLog(r)
	t, _ := template.ParseFiles("html/forum.html")
	t.Execute(w, f)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	Nickname := r.URL.Path[len("/user/"):]
	if !CheckUser(Nickname) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	u := GetUserPage(Nickname)
	u.Own, u.Log = TestLog(r)
	t, _ := template.ParseFiles("html/user.html")
	t.Execute(w, u)
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	Name := r.URL.Path[len("/group/"):]
	if !CheckGroup(Name) {
		http.Redirect(w, r, "/index/", http.StatusFound)
	}
	g := GetGroupPage(Name)
	g.Own, g.Log = TestLog(r)
	t, _ := template.ParseFiles("html/group.html")
	t.Execute(w, g)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	l := GetLoginPage()
	l.Own, l.Log = TestLog(r)
	t, _ := template.ParseFiles("html/login.html")
	t.Execute(w, l)
}

func LoginCookie(Nickname string) http.Cookie {
	cookieValue := Nickname + ":" + GetMd5(Nickname+strconv.Itoa(rand.Intn(100000000)))
	expire := time.Now().AddDate(0, 0, 1)
	return http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire, HttpOnly: true, Path: "/"}
}

func TestLog(r *http.Request) (*PrUser, bool) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		return nil, false
	}
	z := strings.Split(cookie.Value, ":")
	nickname := z[0]
	if !CheckUser(nickname) {
		return nil, false
	}
	user := GetUser(nickname)
	if user.Cookie != cookie.Value {
		return nil, false
	}
	return &PrUser{user, GetGroup(user.GroupName)}, true
}

func EnterHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")
	pass := r.FormValue("password")
	if !CheckUser(nickname) {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}
	user := GetUser(nickname)
	s := GetMd5(pass)
	if user.Hash != s {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}
	cookie := LoginCookie(nickname)
	http.SetCookie(w, &cookie)
	user.Cookie = cookie.Value
	fmt.Println("Loaded COOKIES: " + cookie.Value)
	user.Update()
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func NewgroupHandler(w http.ResponseWriter, r *http.Request) {
	newgroup := GetNewgroupPage()
	newgroup.Own, newgroup.Log = TestLog(r)
	t, _ := template.ParseFiles("html/newgroup.html")
	t.Execute(w, newgroup)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	reg := GetRegisterPage()
	reg.Own, reg.Log = TestLog(r)
	t, _ := template.ParseFiles("html/register.html")
	t.Execute(w, reg)
}

func DoneregHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")
	pass := r.FormValue("password")
	reppass := r.FormValue("reppassword")
	if nickname == "" || pass == "" {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	if CheckUser(nickname) || pass != reppass {
		http.Redirect(w, r, "/register/", http.StatusFound)
		return
	}
	info := r.FormValue("info")
	cookie := LoginCookie(nickname)
	user := User{nickname, info, GetMd5(pass), cookie.Value, "DEFAULT"}
	user.Insert()
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func CreategroupHandler(w http.ResponseWriter, r *http.Request) {
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	name := r.FormValue("name")
	admin := r.FormValue("admin") == "True"
	untouchable := r.FormValue("untouchable") == "True"
	readonly := r.FormValue("readonly") == "True"
	moderator := r.FormValue("moderator") == "True"
	if name == "" {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	group := Group{name, untouchable, moderator, readonly, admin}
	group.Insert()
	http.Redirect(w, r, "/groups/", http.StatusFound)
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/send/"):]
	ThemeId, err := strconv.Atoi(StringId)
	if err != nil || !CheckTheme(ThemeId) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	prUz, ok := TestLog(r)
	if !ok || prUz.Group.ReadOnly {
		http.Redirect(w, r, "/theme/"+strconv.Itoa(ThemeId), http.StatusFound)
	}
	text := r.FormValue("message")
	id := 1
	m := Message{id, text, time.Now().Format(time.RFC3339), prUz.User.Nickname, ThemeId}
	m.Insert()
	http.Redirect(w, r, "/theme/"+strconv.Itoa(ThemeId), http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "SessionID",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func NewthemeHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/newtheme/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckSubforum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	newtheme := GetNewthemePage(Id)
	newtheme.Own, newtheme.Log = TestLog(r)
	if !newtheme.Log || newtheme.Own.Group.ReadOnly {
		http.Redirect(w, r, "/index/", http.StatusFound)
	}
	t, _ := template.ParseFiles("html/newtheme.html")
	t.Execute(w, newtheme)
}

func CreatethemeHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/createtheme/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckSubforum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	prUs, _ := TestLog(r)
	name := r.FormValue("name")
	description := r.FormValue("description")
	theme := Theme{1, name, description, prUs.User.Nickname, Id}
	theme.Insert()
	http.Redirect(w, r, "/subforum/"+strconv.Itoa(Id), http.StatusFound)
}

func UsersaveHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Path[len("/saveuser/"):]
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin || !CheckUser(nickname) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	group := r.FormValue("group")
	if group == "" {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	info := r.FormValue("info")
	user := GetUser(nickname)
	user.GroupName = group
	user.Info = info
	user.Update()
	http.Redirect(w, r, "/user/"+nickname, http.StatusFound)
}

func SaveeditmessageHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/saveeditmessage/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckMessage(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	Own, Log := TestLog(r)
	if err != nil || !Log || (!Own.Group.Admin && !Own.Group.Moderator) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	text := r.FormValue("message")
	if text == "" {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	message := GetMessage(Id)
	message.Text = text
	message.Update()
	http.Redirect(w, r, "/theme/"+strconv.Itoa(message.ThemeId), http.StatusFound)
}

func NewsubforumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/newsubforum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckForum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	newsubforum := GetNewsubforumPage(Id)
	newsubforum.Own, newsubforum.Log = TestLog(r)
	if !newsubforum.Log || !newsubforum.Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
	}
	t, _ := template.ParseFiles("html/newsubforum.html")
	t.Execute(w, newsubforum)
}

func CreatesubforumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/createsubforum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckForum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	name := r.FormValue("name")
	description := r.FormValue("description")
	subforum := Subforum{1, name, description, Id}
	subforum.Insert()
	http.Redirect(w, r, "/forum/"+strconv.Itoa(Id), http.StatusFound)
}

func NewforumHandler(w http.ResponseWriter, r *http.Request) {
	newforum := GetNewforumPage()
	newforum.Own, newforum.Log = TestLog(r)
	if !newforum.Log || !newforum.Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("html/newforum.html")
	t.Execute(w, newforum)
}

func CreateforumHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	forum := Forum{1, name}
	forum.Insert()
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func DeleteforumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/deleteforum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckForum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	DeleteForum(Id)
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func DeletesubforumHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/deletesubforum/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckSubforum(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	DeleteSubforum(Id)
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	members := GetMembersPage()
	members.Own, members.Log = TestLog(r)
	t, _ := template.ParseFiles("html/members.html")
	t.Execute(w, members)
}

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	groups := GetAllgroupsPage()
	groups.Own, groups.Log = TestLog(r)
	t, _ := template.ParseFiles("html/allgroups.html")
	t.Execute(w, groups)
}

func DeleteuserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Path[len("/deleteuser/"):]
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin || !CheckUser(nickname) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	DeleteUser(nickname)
	http.Redirect(w, r, "/members/", http.StatusFound)
}

func DeletegroupHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/deletegroup/"):]
	Own, Log := TestLog(r)
	if !Log || !Own.Group.Admin || !CheckGroup(name) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	DeleteGroup(name)
	http.Redirect(w, r, "/groups/", http.StatusFound)
}

func UsereditHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Path[len("/edituser/"):]
	if !CheckUser(nickname) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	user := GetUsereditPage(nickname)
	user.Own, user.Log = TestLog(r)
	if !user.Log || !user.Own.Group.Admin {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("html/edituser.html")
	t.Execute(w, user)
}

func DeletemessageHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/deletemessage/"):]
	Id, err := strconv.Atoi(StringId)
	if err != nil || !CheckMessage(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	Own, Log := TestLog(r)
	themeId := GetMessage(Id).ThemeId
	if !Log || (!Own.Group.Admin && !Own.Group.Moderator) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	DeleteMessage(Id)
	http.Redirect(w, r, "/theme/"+strconv.Itoa(themeId), http.StatusFound)
}

func EditmessageHandler(w http.ResponseWriter, r *http.Request) {
	StringId := r.URL.Path[len("/editmessage/"):]
	Id, err := strconv.Atoi(StringId)
	Own, Log := TestLog(r)
	if err != nil || !Log || (!Own.Group.Admin && !Own.Group.Moderator) || !CheckMessage(Id) {
		http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
	page := GetEditmessagePage(Id)
	t, _ := template.ParseFiles("html/editmessage.html")
	t.Execute(w, page)
}

func main() {
	http.HandleFunc("/saveeditmessage/", SaveeditmessageHandler)
	http.HandleFunc("/editmessage/", EditmessageHandler)
	http.HandleFunc("/deletemessage/", DeletemessageHandler)
	http.HandleFunc("/saveuser/", UsersaveHandler)
	http.HandleFunc("/edituser/", UsereditHandler)
	http.HandleFunc("/deletegroup/", DeletegroupHandler)
	http.HandleFunc("/creategroup/", CreategroupHandler)
	http.HandleFunc("/groups/", GroupsHandler)
	http.HandleFunc("/newgroup/", NewgroupHandler)
	http.HandleFunc("/deleteuser/", DeleteuserHandler)
	http.HandleFunc("/members/", MembersHandler)
	http.HandleFunc("/deletesubforum/", DeletesubforumHandler)
	http.HandleFunc("/deleteforum/", DeleteforumHandler)
	http.HandleFunc("/createsubforum/", CreatesubforumHandler)
	http.HandleFunc("/newsubforum/", NewsubforumHandler)
	http.HandleFunc("/createforum/", CreateforumHandler)
	http.HandleFunc("/newforum/", NewforumHandler)
	http.HandleFunc("/createtheme/", CreatethemeHandler)
	http.HandleFunc("/newtheme/", NewthemeHandler)
	http.HandleFunc("/logout/", LogoutHandler)
	http.HandleFunc("/send/", SendHandler)
	http.HandleFunc("/donereg/", DoneregHandler)
	http.HandleFunc("/register/", RegisterHandler)
	http.HandleFunc("/enter/", EnterHandler)
	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/group/", GroupHandler)
	http.HandleFunc("/user/", UserHandler)
	http.HandleFunc("/index/", IndexHandler)
	http.HandleFunc("/subforum/", SubforumHandler)
	http.HandleFunc("/theme/", ThemeHandler)
	http.HandleFunc("/forum/", ForumHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(":8090", nil)
}
