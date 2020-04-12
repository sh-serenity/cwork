package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"os/exec"
	"os"
	"net/http"
)
func regprocHandle2(w http.ResponseWriter, r *http.Request) {

	var user User
    	user = isauth(w,r)
	fmt.Println(user.id);
	chknon(w,r)
	db := dbConnect()
	r.ParseForm()
	var rchk regchk
	var resume string
	var note_username, note_usernamep, note_password, note_invite, note_fname, note_sname string
	var username = r.FormValue("login")
	var password = r.FormValue("password")
	var fname = r.FormValue("fname")
	var sname = r.FormValue("sname")
	var conifirm = r.FormValue("conifirm") 
	if user.id == user.aid  {
		note_invite = "Права есть"
		rchk.Invitech = 1
	} else {
		note_invite = "Прав нет"
		rchk.Invitech = 0
	}
	eu := validlogin.FindStringSubmatch(username)
	if eu == nil {
		note_username = "Юзернейм неверный"
		rchk.Usernameex = 0
	} else {
		note_username = "Юзернейм в порядке"
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
		note_usernamep = "Такой логин еще не зарегистрирован"
		rchk.Usernamereg = 1
	} else {
		note_usernamep = "Логин занят"
		rchk.Usernamereg = 0
	}

	pu := validStr.FindStringSubmatch(password)
	if pu == nil {
		note_password = "Пароль может содержать только заглавные и маленькие буквы латиницы, и цифры."
		rchk.Passwordrx = 0
	} else {
		note_password = "Пароль содержит правильные символы."
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
	var url, to string
	fn := validpass.FindStringSubmatch(fname)
	if fn == nil {
		note_fname = "Имя может содержать только заглавные и маленькие буквы,."
		rchk.fnrx = 0
	} else {
		note_username = "Имя содержит правильные символы."
		rchk.fnrx = 1
	}
	sn := validpass.FindStringSubmatch(sname)
	if sn == nil {
		note_sname = "Фамилия может содержать только заглавные и маленькие буквы,."
		rchk.snrx = 0
	} else {
		note_sname = "Фамилия содержит правильные символы."
		rchk.snrx = 1
	}

	if len(sname) < 64 {
		rchk.sln = 1
	}else {

		rchk.sln = 0
	}

	if rchk.Usernamereg == 1 && rchk.Usernameex == 1 && rchk.Passwordcon == 1 && rchk.Passwordrx == 1 && rchk.Invitech == 1 && rchk.fnrx == 1 && rchk.snrx == 1 && rchk.Invitech ==  1{

		result, err := db.Exec("insert into users (username,password,fname,sname,rootid) values(?,MD5(?),?,?,?)", username, password, fname, sname,user.aid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result.LastInsertId()) // id добавленного объекта
		fmt.Println(result.RowsAffected())

		otblin := &exec.Cmd{Path:"/usr/bin/sudo",Args:[]string{"/bin/sh","/root/au.sh",username, user.username, password},Stdout:os.Stdout,Stderr:os.Stderr}
		otblin.Run()
		if err != nil {
	    	    log.Fatal(err)
		}

		http.Redirect(w,r,"/home/",301)
	} else {
		resume = note_invite + " " + note_username +" " + note_usernamep + " " + note_fname + " "  + note_sname + " " + note_password + " Данные введены с ошибками. Поправьте и попробуйте снова."
		url = "/reg2/"
		to = "Регистрация"
		p := &regdata{Resume: resume,Url: url, To: to}
		t, _ := template.ParseFiles("tmpl/regproc.html","tmpl/header.html","tmpl/footer.html")
		t.ExecuteTemplate(w,"index", p)

	}

}