package main
import (

		"html/template"
		"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"regexp"

	//	"github.com/shurcooL/github_flavored_markdown"
	//	"html/template"
	"net"
	"net/http"
	"net/http/fcgi"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("17923641793298746918723649781112")
	store = sessions.NewCookieStore(key)
)

type tmp1 struct {
	Tmp string
}

func IndexHand(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("tmpl/body.html","tmpl/header.html", "tmpl/footer.html")
	t.ExecuteTemplate(w,"index", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func EnterHand(w http.ResponseWriter, r *http.Request) {

    t, err := template.ParseFiles("tmpl/enter.html","tmpl/header.html", "tmpl/footer.html")
    t.ExecuteTemplate(w,"index", nil)
    if err != nil {
	fmt.Println(err)
    }
}

func RegForm(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("tmpl/regform.html","tmpl/header.html", "tmpl/footer.html")
	t.ExecuteTemplate(w,"regform", nil)
	if err != nil {
		fmt.Println(err)
	}

}

func RegForm2(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("tmpl/regform2.html","tmpl/header.html", "tmpl/footer.html")
	t.ExecuteTemplate(w,"regform", nil)
	if err != nil {
		fmt.Println(err)
	}

}
var validnon = regexp.MustCompile("^/(reg|regproc|enter|sign|exit|home|regproc2|reg2|reghelp|userhelp)/$")
var vaitdn = regexp.MustCompile("^/(comform|comment|users)/([0-9]+)$")

func chknon(w http.ResponseWriter, r *http.Request)  {
	m := validnon.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w,r,"/static/404.htmml",301)
	}
}

func chkn(w http.ResponseWriter, r *http.Request) {
	m := vaitdn.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w,r,"/static/404.htmml",301)
	}
}

type tmp struct {
	Title, Note  string
}

func leaveHandler(w http.ResponseWriter, r *http.Request) {
	chknon(w,r)
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w,r,"/",301)
}

func RegHelp(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("tmpl/reghelp.html","tmpl/header.html", "tmpl/footer.html")
    if err != nil {
	fmt.Println(err)
    }
    t.ExecuteTemplate(w,"index", nil)
    if err != nil {
	fmt.Println(err)
    }
}

func UserHelp(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("tmpl/userhelp.html","tmpl/header.html", "tmpl/footer.html")
    if err != nil {
	fmt.Println(err)
    }
    t.ExecuteTemplate(w,"index", nil)
    if err != nil {
	fmt.Println(err)
    }
}



func main() {
//	fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/", IndexHand);
	http.HandleFunc("/enter/", EnterHand)
	http.HandleFunc("/reg/", RegForm)
	http.HandleFunc("/reg2/", RegForm2)
	http.HandleFunc("/reghelp/", RegHelp)
	http.HandleFunc("/home/", Home)
	http.HandleFunc("/regproc/", regprocHandle)
	http.HandleFunc("/regproc2/", regprocHandle2)
	http.HandleFunc("/sign/",signHandler)
	http.HandleFunc("/userhelp/",UserHelp)
	http.HandleFunc("/exit/", leaveHandler)
	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		return
	}
	fcgi.Serve(l, nil)

}