package main

import (
//	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"os/exec"
	"os"
	//	"log"
	"net/http"
//	"net/smtp"
)
/*
func mailsend(fname string, sname string,email string, msg string) {
	// Set up authentication information.

	smtpServer := mailserver
	auth := smtp.PlainAuth(
		"",
		mailuser,
		mailpass,
		smtpServer,
	)

	from := mail.Address{"noreplay", "noreplay@serenity-net.org"}
	to := mail.Address{fname + " " + sname, email}
	title := "Запрос на подтверждение email от serenity-net.org"

	body := msg;

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer + ":25",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
		//[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
} */

func regprocHandle2(w http.ResponseWriter, r *http.Request) {
	var user User
	user = isauth(w, r)
	chknon(w,r)
	db := dbConnect()
	r.ParseForm()
	//	var login, password string
	//	var fmchk check;
	//	var i int
	//	i = 0
	var rchk regchk
	var resume string
	var note_username, note_password string
//	var invite = r.FormValue("invite")
	var username = r.FormValue("login")
	var password = r.FormValue("password")
	var fname = r.FormValue("fname")
	var sname = r.FormValue("sname")
	var conifirm = r.FormValue("conifirm")
//	var note_invite string
//	var count int
		//var invch string
/*	fmt.Printf("%s = %s", password, conifirm)
	err := db.QueryRow("select COUNT(*) from invites where invite = ?", invite).Scan(&count)
	fmt.Printf("Number of rows are %d\n", count) *?	if err != nil {
		fmt.Println(err)
	} */
	if user.id == user.rootid  {
		//note_invite = "Правав есть"
		rchk.Invitech = 1
	} else {
		//note_invite = "GHfd ytn"
		rchk.Invitech = 0
	}
//	var validemail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	eu := validStr.FindStringSubmatch(username)
	if eu == nil {
		note_username = "юзернейм неверный"
		rchk.Usernameex = 0
	} else {
		note_username = "юзернейм в порядке"
		rchk.Usernameex = 1
	}
//	rchk.Usernameex = 1
	var ucount int
	ucount = 0
	uerr := db.QueryRow("select COUNT(*) from users where username = ?", username).Scan(&ucount)
	fmt.Printf("Number of rows are %d\n", ucount)
	if uerr != nil {
		fmt.Println(uerr)
	}
	if ucount == 0 {
		note_username = note_username + "Такой логин еще не зарегистрирован"
		rchk.Usernamereg = 1
	} else {
		note_username = "Вы уже регистрировались здесь, с тем же логином"
		rchk.Usernamereg = 0
	}

	pu := validStr.FindStringSubmatch(password)
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
		note_password = "Имя может содержать только заглавные и маленькие буквы,."
		rchk.fnrx = 0
	} else {
		note_username = "Имя содержит правильные символы."
		rchk.fnrx = 1
	}

	if len(fname) < 64 {
		rchk.fln = 1
	} else {

		rchk.fln = 0
	}




	sn := validStr.FindStringSubmatch(sname)

	if sn == nil {
		note_password = "Фамилия может содержать только заглавные и маленькие буквы,."
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
/*		result, err := db.Exec("insert into users (username, password) values(?,MD5(?))", username, password)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result.LastInsertId()) // id добавленного объекта
		fmt.Println(result.RowsAffected()) */
		str,err := GenerateRandomString(40)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(str)
		result, err := db.Exec("insert into users (username,password,fname,sname,rootid values(?,MD5(?),?,?,?)", username, password, fname, sname,user.rootid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result.LastInsertId()) // id добавленного объекта
		fmt.Println(result.RowsAffected())
		//msg  := "<html lang=\"ru_RU\"><head><title>Подтверждение почты от serenity-net.org</title></head><body>Подтвердите вашу регистрацию на сайте serenity-net.org, перейдя по <a href=\"http://serenity-net.org/checkemail/" + str +"\">ссылке</a></body></html>"
		/* msg := "Подтвердите вашу регистрацию на сайте serenity-net.org, перейдя по ссылке http://k8s.serenity-net.org/checkemail/" + str
		if err != nil {
			fmt.Println(err)
		}
		mailsend(fname, sname, username, msg)

		resume = "На ваш email: " + username + " отправенна ссылка для его подтвержедения. Перейдите по ссылке, после этого вы получите возможность заходить на сайт."
		p := &trash{Title: "Завершение регистрации", Note: resume}
		t, _ := template.ParseFiles("tmpl/tmp1.html")
		t.Execute(w, p) */

	otblin := &exec.Cmd{Path:"/usr/bin/sudo",Args:[]string{"/bin/sh","/root/au.sh",username, user.username, password},Stdout:os.Stdout,Stderr:os.Stderr}
	otblin.Run()
	if err != nil {
	    log.Fatal(err)
	}


		//fmt.Printf(out)

		http.Redirect(w,r,"/enter/",301)
	} else {
		resume = "Данные введены с ошибками. Поправьте и попробуйте снова. Или у вас прав нету."
		url = "/reg2/"
		to = "Регистрация"
		//	end = 0
		p := &regdata{Resume: resume,Url: url, To: to}
		t, _ := template.ParseFiles("tmpl/regproc.html","tmpl/header.html","tmpl/footer.html")
		t.ExecuteTemplate(w,"index", p)

	}

}