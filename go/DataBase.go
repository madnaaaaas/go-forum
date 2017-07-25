package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var dbuser, dbpassword, dbname string = "root", "qweasd", "forumtemp"

type User struct {
	Nickname  string
	Info      string
	Hash      string
	Cookie    string
	GroupName string
}

func (p User) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _User values (?,?,?,?,?)", p.Nickname, p.Info, p.Hash, p.Cookie, p.GroupName)
	db.Close()
}

func (p User) Update() {
	if p.Nickname == "DELETED" {
		return
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _User set Info=? , GroupName=?, Hash=?, Cookie=? where NickName=?", p.Info, p.GroupName, p.Hash, p.Cookie, p.Nickname)
	db.Close()
}

func (p User) Print() {
	fmt.Printf("%s %s %s %s\n", p.Nickname, p.Info, p.Hash, p.GroupName)
}

func CheckUser(Nickname string) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _User where Nickname=?", Nickname).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetUser(Nickname string) User {
	var u User
	u.Nickname = Nickname
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select Info,Hash,GroupName,Cookie from _User where Nickname=?", Nickname).Scan(&(u.Info), &(u.Hash), &(u.GroupName), &(u.Cookie))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such User\n")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		//u.Print()
	}
	db.Close()
	return u
}

func GetAllUsers() []User {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _User")
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&(u.Nickname), &(u.Info), &(u.Hash), &(u.Cookie), &(u.GroupName))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, u)
	}
	rows.Close()
	db.Close()
	return mas
}

func GetAllUsersInGroup(Name string) []User {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _User where GroupName=?", Name)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&(u.Nickname), &(u.Info), &(u.Hash), &(u.Cookie), &(u.GroupName))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, u)
	}
	rows.Close()
	db.Close()
	return mas
}

func DeleteUser(Nickname string) {
	group := GetGroup(GetUser(Nickname).GroupName)
	if group.Untouchable || Nickname == "DELETED" {
		return
	}
	themes := GetAllThemesFromUser(Nickname)
	for _, v := range themes {
		v.CreatorNickname = "DELETED"
		v.Update()
	}
	messages := GetAllMessagesFromUser(Nickname)
	for _, v := range messages {
		v.CreatorNickname = "DELETED"
		v.Update()
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _User where Nickname=?", Nickname)
	db.Close()
}

type Group struct {
	Name        string
	Untouchable bool
	Moderator   bool
	ReadOnly    bool
	Admin       bool
}

func (p Group) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _Group values (?,?,?,?,?)", p.Name, p.Untouchable, p.Moderator, p.ReadOnly, p.Admin)
	db.Close()
}

func (p Group) Update() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _Group set Untouchable=? , Moderator=?, ReadOnly=?, Admin=? where Name=?", p.Untouchable, p.Moderator, p.ReadOnly, p.Admin, p.Name)
	db.Close()
}

func (p Group) Print() {
	fmt.Printf("%s %b %b %b %b\n", p.Name, p.Untouchable, p.Moderator, p.ReadOnly, p.Admin)
}

func CheckGroup(Name string) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _Group where Name=?", Name).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetGroup(Name string) Group {
	var g Group
	g.Name = Name
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select Untouchable,Moderator,ReadOnly,Admin from _Group where Name=?", Name).Scan(&(g.Untouchable), &(g.Moderator), &(g.ReadOnly), &(g.Admin))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such Group\n")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		//g.Print()
	}
	db.Close()
	return g
}

func GetAllGroups() []Group {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Group")
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Group, 0)
	for rows.Next() {
		var g Group
		err := rows.Scan(&(g.Name), &(g.Untouchable), &(g.Moderator), &(g.ReadOnly), &(g.Admin))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, g)
	}
	rows.Close()
	db.Close()
	return mas
}

func DeleteGroup(Name string) {
	if Name == "DEFAULT" || Name == "DELETED" || Name == "Admin" {
		return
	}
	users := GetAllUsersInGroup(Name)
	for _, v := range users {
		v.GroupName = "DEFAULT"
		v.Update()
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _Group where Name=?", Name)
	db.Close()
}

type Forum struct {
	Id   int
	Name string
}

func (p Forum) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _Forum (Name) values (?)", p.Name)
	db.Close()
}

func (p Forum) Update() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _Forum set Name=? where Id=?", p.Name, p.Id)
	db.Close()
}

func (p Forum) Print() {
	fmt.Printf("%d %s\n", p.Id, p.Name)
}

func CheckForum(Id int) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _Forum where Id=?", Id).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetForum(Id int) Forum {
	var f Forum
	f.Id = Id
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select Name from _Forum where Id=?", Id).Scan(&(f.Name))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such Forum\n")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		//f.Print()
	}
	db.Close()
	return f
}

func GetAllForums() []Forum {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Forum")
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Forum, 0)
	for rows.Next() {
		var f Forum
		err := rows.Scan(&(f.Id), &(f.Name))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, f)
	}
	rows.Close()
	db.Close()
	return mas
}

func DeleteForum(Id int) {
	mas := GetAllSubforums(Id)
	for _, v := range mas {
		DeleteSubforum(v.Id)
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _Forum where Id=?", Id)
	db.Close()
}

type Theme struct {
	Id              int
	Name            string
	Description     string
	CreatorNickname string
	SubforumId      int
}

func (p Theme) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _Theme (Name,Description,Creator,SubforumId) values (?,?,?,?)", p.Name, p.Description, p.CreatorNickname, p.SubforumId)
	db.Close()
}

func (p Theme) Update() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _Theme set Name=?, Description=?, Creator=?, SubforumId=? where Id=?", p.Name, p.Description, p.CreatorNickname, p.SubforumId, p.Id)
	db.Close()
}

func (p Theme) Print() {
	fmt.Printf("%d %s %s %s %d\n", p.Id, p.Name, p.Description, p.CreatorNickname, p.SubforumId)
}

func CheckTheme(Id int) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _Theme where Id=?", Id).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetTheme(Id int) Theme {
	var t Theme
	t.Id = Id
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select Name,Description,Creator,SubforumId from _Theme where Id=?", Id).Scan(
		&(t.Name), &(t.Description), &(t.CreatorNickname), &(t.SubforumId))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such Theme\n")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		//t.Print()
	}
	db.Close()
	return t
}

func GetAllThemes(SubforumId int) []Theme {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Theme where SubforumId=?", SubforumId)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Theme, 0)
	for rows.Next() {
		var t Theme
		err := rows.Scan(&(t.Id), &(t.Name), &(t.Description), &(t.CreatorNickname), &(t.SubforumId))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, t)
	}
	rows.Close()
	db.Close()
	return mas
}

func DeleteTheme(Id int) {
	mas := GetAllMessages(Id)
	for _, v := range mas {
		DeleteMessage(v.Id)
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _Theme where Id=?", Id)
	db.Close()
}

func GetAllThemesFromUser(Nickname string) []Theme {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Theme where Creator=?", Nickname)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Theme, 0)
	for rows.Next() {
		var t Theme
		err := rows.Scan(&(t.Id), &(t.Name), &(t.Description), &(t.CreatorNickname), &(t.SubforumId))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, t)
	}
	rows.Close()
	db.Close()
	return mas
}

type Subforum struct {
	Id          int
	Name        string
	Description string
	ForumId     int
}

func (p Subforum) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _Subforum (Name,Description,ForumId) values (?,?,?)", p.Name, p.Description, p.ForumId)
	db.Close()
}

func (p Subforum) Update() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _Subforum set Name=?, Description=?, ForumId=? where Id=?", p.Name, p.Description, p.ForumId, p.Id)
	db.Close()
}

func (p Subforum) Print() {
	fmt.Printf("%d %s %s %d\n", p.Id, p.Name, p.Description, p.ForumId)
}

func CheckSubforum(Id int) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _Subforum where Id=?", Id).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetSubforum(Id int) Subforum {
	var s Subforum
	s.Id = Id
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select Name,Description,ForumId from _Subforum where Id=?", Id).Scan(&(s.Name), &(s.Description), &(s.ForumId))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such Subforum\n")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		//s.Print()
	}
	db.Close()
	return s
}

func GetAllSubforums(ForumId int) []Subforum {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Subforum where ForumId=?", ForumId)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Subforum, 0)
	for rows.Next() {
		var s Subforum
		err := rows.Scan(&(s.Id), &(s.Name), &(s.Description), &(s.ForumId))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, s)
	}
	rows.Close()
	db.Close()
	return mas
}

func DeleteSubforum(Id int) {
	mas := GetAllThemes(Id)
	for _, v := range mas {
		DeleteTheme(v.Id)
	}
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _Subforum where Id=?", Id)
	db.Close()
}

type Message struct {
	Id              int
	Text            string
	Date            string
	CreatorNickname string
	ThemeId         int
}

func (p Message) Insert() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("insert into _Message (_Text,_Date,Creator,ThemeId) values (?,?,?,?)", p.Text, p.Date, p.CreatorNickname, p.ThemeId)
	db.Close()
}

func (p Message) Update() {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("update _Message set _Text=?, _Date=?, Creator=?, ThemeId=? where Id=?", p.Text, p.Date, p.CreatorNickname, p.ThemeId, p.Id)
	db.Close()
}

func (p Message) Print() {
	fmt.Printf("%d %s %s %s %d\n", p.Id, p.Text, p.Date, p.CreatorNickname, p.ThemeId)
}

func CheckMessage(Id int) bool {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	num := 0
	err := db.QueryRow("select count(*) from _Message where id=?", Id).Scan(&num)
	db.Close()
	return err == nil && num != 0
}

func GetMessage(Id int) Message {
	var m Message
	m.Id = Id
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	err := db.QueryRow("select _Text,_Date,Creator,ThemeId from _Message where Id=?", Id).Scan(&(m.Text), &(m.Date), &(m.CreatorNickname), &(m.ThemeId))
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No Such Message\n")
		break
	case err != nil:
		log.Fatal(err)
	}
	return m
}

func DeleteMessage(Id int) {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	db.Exec("delete from _Message where Id=?", Id)
	db.Close()
}

func GetAllMessages(ThemeId int) []Message {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Message where ThemeId=?", ThemeId)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Message, 0)
	for rows.Next() {
		var m Message
		err := rows.Scan(&(m.Id), &(m.Text), &(m.Date), &(m.CreatorNickname), &(m.ThemeId))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, m)
	}
	rows.Close()
	db.Close()
	return mas
}

func GetAllMessagesFromUser(Nickname string) []Message {
	db, _ := sql.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname)
	rows, err := db.Query("select * from _Message where Creator=?", Nickname)
	if err != nil {
		log.Fatal(err)
	}
	mas := make([]Message, 0)
	for rows.Next() {
		var m Message
		err := rows.Scan(&(m.Id), &(m.Text), &(m.Date), &(m.CreatorNickname), &(m.ThemeId))
		if err != nil {
			log.Fatal(err)
		}
		mas = append(mas, m)
	}
	rows.Close()
	db.Close()
	return mas
}
