package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/robster0/crud/sqldb"
)

type Template struct {
	Title       string
	Description string
	Posts       []sqldb.Post
	Comments    []sqldb.Comment
	ID          int
	N           string
	C           string
	Created     string
	V           string
	Text        string
}

func Home_GET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page")

	executeTemplate(Template{Title: "Home"}, "template/home.gohtml", w)
}
func CreateP_GET(w http.ResponseWriter, r *http.Request) {

	executeTemplate(Template{V: "p", Text: "post"}, "template/create.gohtml", w)
}
func CreateC_GET(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)["id"]

	executeTemplate(Template{V: fmt.Sprint("c/", params), Text: "comment"}, "template/create.gohtml", w)
}
func Read_GET(w http.ResponseWriter, r *http.Request) {

	posts, err := sqldb.SelectAll()

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return
	}

	comments, err := sqldb.GetComments()

	fmt.Println("COmments: ")
	fmt.Println(comments)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return
	}

	executeTemplate(Template{Posts: posts, Comments: comments}, "template/read.gohtml", w)
}
func ReadOne_GET(w http.ResponseWriter, r *http.Request) {

	//params := mux.Vars(r)["id"]

	executeTemplate(Template{Title: "Read Page"}, "template/read.gohtml", w)
}
func Update_GET(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)["id"]

	id, err := strconv.Atoi(params)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
	}

	posts, err := sqldb.SelectOne(id)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
	}

	executeTemplate(Template{ID: posts.ID, N: posts.N, C: posts.C}, "template/update.gohtml", w)
}

func Delete_GET(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)["id"]

	p, err := strconv.Atoi(params)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
	}

	c := sqldb.Delete(p)

	if !c {
		http.Redirect(w, r, "/error?q=there seem to be an issue right now", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/read", http.StatusSeeOther)
	//executeTemplate(Template{Title: "Delete page"}, "template/delete.gohtml", w)
}

func Error_GET(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")

	xss := regexp.MustCompile(`(<script>|<\/script>)`)

	if xss.MatchString(q) {
		http.Redirect(w, r, "/error?q=no xss in params", http.StatusSeeOther)
	}

	var title string
	var description string

	if q == "" {
		title = "Page not found"
		description = "the page you requested does not exist"
	} else {
		title = "Error"
		description = q
	}

	executeTemplate(Template{Title: title, Description: description}, "template/error.gohtml", w)
}

func CreateP_POST(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	err := sqldb.Insert(r.Form["Name"][0], r.Form["Content"][0], 0, true)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func CreateC_POST(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	params := mux.Vars(r)["id"]

	n, err := strconv.Atoi(params)

	fmt.Println(n)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return
	}

	err = sqldb.Insert(r.Form["Name"][0], r.Form["Content"][0], n, false)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/read", http.StatusSeeOther)

}
func Update_POST(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	params := mux.Vars(r)["id"]

	id, err := strconv.Atoi(params)

	if err != nil {
		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
	}

	if err = sqldb.Update(r.Form["Name"][0], r.Form["Content"][0], id); err != nil {

		http.Redirect(w, r, fmt.Sprint("/error?q=", err), http.StatusSeeOther)
		return

	}

	http.Redirect(w, r, "/read", http.StatusSeeOther)
}
func executeTemplate(data Template, path string, w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles(path))

	if err := tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
	}
}
