package main

import ()

type PrUser struct {
	User  User
	Group Group
}

type ForumPage struct {
	Title Forum
	Body  []Subforum
	Own   *PrUser
	Log   bool
}

func GetForumPage(Id int) ForumPage {
	mas := GetAllSubforums(Id)
	return ForumPage{GetForum(Id), mas, nil, false}
}

type IndexPage struct {
	Title string
	Body  []ForumPage
	Own   *PrUser
	Log   bool
}

func GetIndexPage() IndexPage {
	fms := GetAllForums()
	b := make([]ForumPage, 0)
	for _, v := range fms {
		mas := GetAllSubforums(v.Id)
		b = append(b, ForumPage{v, mas, nil, false})
	}
	return IndexPage{"Index", b, nil, false}
}

type SubforumPage struct {
	Title Subforum
	Forum Forum
	Body  []Theme
	Own   *PrUser
	Log   bool
}

func GetSubforumPage(Id int) SubforumPage {
	mas := GetAllThemes(Id)
	this := GetSubforum(Id)
	return SubforumPage{this, GetForum(this.ForumId), mas, nil, false}
}

type ThemePage struct {
	Title    Theme
	Subforum Subforum
	Forum    Forum
	Body     []Message
	Own      *PrUser
	Log      bool
}

func GetThemePage(Id int) ThemePage {
	mas := GetAllMessages(Id)
	this := GetTheme(Id)
	sub := GetSubforum(this.SubforumId)
	fo := GetForum(sub.ForumId)
	return ThemePage{this, sub, fo, mas, nil, false}
}

type EditmessagePage struct {
	ThemePage
	Edited Message
}

func GetEditmessagePage(Id int) EditmessagePage {
	mes := GetMessage(Id)
	return EditmessagePage{GetThemePage(mes.ThemeId), mes}
}

type UserPage struct {
	Title string
	Body  User
	Own   *PrUser
	Log   bool
}

func GetUserPage(Nickname string) UserPage {
	return UserPage{Nickname, GetUser(Nickname), nil, false}
}

type UsereditPage struct {
	UserPage
	Groups []Group
}

func GetUsereditPage(Nickname string) UsereditPage {
	return UsereditPage{GetUserPage(Nickname), GetAllGroups()}
}

type GroupPage struct {
	Title string
	Body  Group
	Own   *PrUser
	Log   bool
}

func GetGroupPage(Name string) GroupPage {
	return GroupPage{Name, GetGroup(Name), nil, false}
}

type LoginPage struct {
	Title string
	Own   *PrUser
	Log   bool
}

func GetLoginPage() LoginPage {
	return LoginPage{"Login", nil, false}
}

type RegisterPage struct {
	Title string
	Own   *PrUser
	Log   bool
}

func GetRegisterPage() RegisterPage {
	return RegisterPage{"Register", nil, false}
}

type NewgroupPage struct {
	Title string
	Own   *PrUser
	Log   bool
}

func GetNewgroupPage() NewgroupPage {
	return NewgroupPage{"New Group", nil, false}
}

type EnterPage struct {
	Title string
	Body  User
	Own   *PrUser
	Log   bool
}

func GetEnterPage(user User) EnterPage {
	return EnterPage{"Welcom", user, nil, false}
}

type NewthemePage struct {
	Title string
	Body  Subforum
	Own   *PrUser
	Log   bool
}

func GetNewthemePage(Id int) NewthemePage {
	return NewthemePage{"NewTheme", GetSubforum(Id), nil, false}
}

type NewforumPage struct {
	Title string
	Own   *PrUser
	Log   bool
}

func GetNewforumPage() NewforumPage {
	return NewforumPage{"NewForum", nil, false}
}

type NewsubforumPage struct {
	Title string
	Body  Forum
	Own   *PrUser
	Log   bool
}

func GetNewsubforumPage(Id int) NewsubforumPage {
	return NewsubforumPage{"NewSubforum", GetForum(Id), nil, false}
}

type MembersPage struct {
	Body []User
	Own  *PrUser
	Log  bool
}

func GetMembersPage() MembersPage {
	return MembersPage{GetAllUsers(), nil, false}
}

type AllgroupsPage struct {
	Body []Group
	Own  *PrUser
	Log  bool
}

func GetAllgroupsPage() AllgroupsPage {
	return AllgroupsPage{GetAllGroups(), nil, false}
}
