package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"os/exec"

	"net/http"
	"regexp"
	"os"
)
type trash struct {
	Title, Note string
}

type regdata struct {
	Invite string
	Username string
	Password string
	Resume string
	Url string
	To string
}

type regchk struct {
	Invitech int
	Usernamereg int
	Usernameex int
	Passwordrx int
	Passwordcon int
	fnrx, snrx int
	pln, fln, sln int
}

var validStr = regexp.MustCompile("^([a-zA-Z0-9а-яА-Я]+)$")
var validlogin = regexp.MustCompile("^([a-z0-9]+)$")
var validpass = regexp.MustCompile("^([a-zA-Z0-9]+)$")
func regprocHandle(w http.ResponseWriter, r *http.Request) {
	chknon(w,r)
	db := dbConnect()
	r.ParseForm()
	var rchk regchk
	var resume string
	var note_username, note_usernamep, note_password, note_invite, note_fname, note_sname string
	var invite = r.FormValue("invite")
	var username = r.FormValue("login")
	var password = r.FormValue("password")
	var fname = r.FormValue("fname")
	var sname = r.FormValue("sname")
	var conifirm = r.FormValue("conifirm")
	var count int
	fmt.Printf("%s = %s", password, conifirm)
	err := db.QueryRow("select COUNT(*) from invites where invite = ?", invite).Scan(&count)
	fmt.Printf("Number of rows are %d\n", count)
	if err != nil {
		fmt.Println(err)
	}
	if count > 0 {
		note_invite = "Инвайт найден"
		rchk.Invitech = 1
	} else {
		note_invite = "Инвайт не найден"
		rchk.Invitech = 0
	}

	eu := validlogin.FindStringSubmatch(username)
	if eu == nil {
		note_username = "юзернейм неверный"
		rchk.Usernameex = 0
	} else {
		note_username = "юзернейм в порядке"
		rchk.Usernameex = 1
	}
	var ucount int
	ucount = 0
	uerr := db.QueryRow("select COUNT(*) from users where username = ?", username).Scan(&ucount)
	fmt.Printf("Number of rows are %d\n", ucount)
	if uerr != nil {
		fmt.Println(uerr)
	}
	if ucount == 0 {
		note_usernamep = note_usernamep + "Такой логин еще не зарегистрирован"
		rchk.Usernamereg = 1
	} else {
		note_usernamep = "Вы уже регистрировались здесь, с тем же логином"
		rchk.Usernamereg = 0
	}

	pu := validpass.FindStringSubmatch(password)
	if pu == nil {
		note_password = "Пароль может содержать только заглавные и маленькие буквы латиницы, и цифры."
		rchk.Passwordrx = 0
	} else {
		note_username = "Пароль содержит правильные символы."
		rchk.Passwordrx = 1
	}

	if password == conifirm {
		note_password = note_password + "Пароль и подтверждение его совападают"
		rchk.Passwordcon = 1
	} else {
		note_password = note_password + "Пароль и подтверждение его не совападают"
		rchk.Passwordcon = 0
	}
	if len(password) < 64 {
	   rchk.pln = 1
	} else
	{
		rchk.pln = 0
	}
	//  var end int;
	var url, to string
	fn := validStr.FindStringSubmatch(fname)
	if fn == nil {
		note_fname = "Имя может содержать только заглавные и маленькие буквы,."
		rchk.fnrx = 0
	} else {
		note_sname = "Имя содержит правильные символы."
		rchk.fnrx = 1
	}

	if len(fname) < 64 {
		rchk.fln = 1
	} else {

		rchk.fln = 0
	}
	sn := validStr.FindStringSubmatch(sname)

	if sn == nil {
		note_sname = "Фамилия может содержать только заглавные и маленькие буквы,."
		rchk.snrx = 0
	} else {
		note_username = "Фамилия содержит правильные символы."
		rchk.snrx = 1
	}

	if len(sname) < 64 {
		rchk.sln = 1
	}else {

		rchk.sln = 0
	}

	if rchk.Usernamereg == 1 && rchk.Usernameex == 1 && rchk.Passwordcon == 1 && rchk.Passwordrx == 1 && rchk.Invitech == 1 && rchk.fnrx == 1 && rchk.snrx == 1 && rchk.Invitech ==  1{
		result, err := db.Exec("insert into users (username,password,fname,sname) values(?,MD5(?),?,?)", username, password, fname, sname)
		rootid,err := result.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result.LastInsertId()) // id добавленного объекта
		fmt.Println(result.RowsAffected())
		result, err = db.Exec("UPDATE users set rootid=? where id=?",rootid,rootid)
		if err != nil {
			fmt.Println(err)
		}

		
		otblin := &exec.Cmd{Path:"/usr/bin/sudo",Args:[]string{"/bin/sh","/opt/go/cwork/au.sh",username, username, password},Stdout:os.Stdout,Stderr:os.Stderr}
		otblin.Run()
		if err != nil {
		    log.Fatal(err)
		}

		http.Redirect(w,r,"/enter/",301)
	} else {
		resume = note_invite + " " + note_username + " " + note_usernamep + " " + note_fname + " "  + note_sname + " " + note_password + " Данные введены с ошибками. Поправьте и попробуйте снова."
		url = "/reg/"
		to = "Регистрация"
		p := &regdata{Resume: resume,Url: url, To: to}
		t, err := template.ParseFiles("tmpl/regproc.html","tmpl/header.html","tmpl/footer.html")
	        if err != nil {
	            fmt.Println(err)
		}
		t.ExecuteTemplate(w,"index", p)
	}

}