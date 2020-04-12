package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"html/template"
	"net/http"
)
type User struct {
	id int
	username string
	fname string
	sname string
	password string
	email string
	about string
	userpic string
	timereg string
	aid int
}


func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w,r,"/",301)
}

func signinname(username string) (User) {
	db := dbConnect()
	var user User
	row := db.QueryRow("select id, username, fname, rootid from users where username = ?", username)
	err := row.Scan(&user.id, &user.username, &user.fname,&user.sname, &user.aid)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return user
}
func signbyid(i int) (User) {
	db := dbConnect()
	var user User
	row := db.QueryRow("select id, username, fname, sname, rootid from users where id = ?", i)
	err := row.Scan(&user.id, &user.username, &user.fname,&user.sname, &user.aid)

	if err != nil {
		fmt.Println(err)}
	db.Close()
	fmt.Println("rootid:",user.aid);
	return user
}
type link struct {
	Title string
	Url string
}


func signHandler(w http.ResponseWriter, r *http.Request) {
	chknon(w,r)
	session, _ := store.Get(r, "cookie-name")

	session.Values["authenticated"] =  false
	db := dbConnect()
	session.Save(r, w)
	r.ParseForm()
	var user User
	var username = r.FormValue("login")
	var password = r.FormValue("password")
	row := db.QueryRow("select id, username, fname, sname, rootid from users where username = ? and password=MD5(?)", username, password)
	err := row.Scan(&user.id, &user.username, &user.fname, &user.sname, &user.aid)
	if err == nil {
		fmt.Println(err)
	}
	if   user.id != 0  {
		 session.Values["userid"] = user.id
		 session.Values["authenticated"] = true
		 session.Save(r, w)
		 http.Redirect(w,r,"/home",301)
	 }

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		tmp := &tmp{Title: "Ошибка",Note:"Или произошла страшная ошибка, или пользователя с такими данными не существует."}
		t, _ := template.ParseFiles("tmpl/tmp.html","tmpl/header.html","tmpl/footer.html")
		t.ExecuteTemplate(w, "tmp",tmp)

	 }
	db.Close()
}

func isauth(w http.ResponseWriter, r *http.Request) (User) {
	session, _ := store.Get(r, "cookie-name")
	var user User
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w,r,"/",302)
	} else {
		userid := session.Values["userid"].(int)
		user = signbyid(userid)
	}
	return user
}

func Home(w http.ResponseWriter, r *http.Request) {

	var user User
	user = isauth(w, r)
	p := &tmp1{Tmp: user.fname}
	t, err := template.ParseFiles("tmpl/home.html","tmpl/header.html", "tmpl/menu.html" ,"tmpl/footer.html")
	if err != nil {
		fmt.Println(err)
	}
	t.ExecuteTemplate(w,"index",p)
	if err != nil {
		fmt.Println(err)
	}


}
