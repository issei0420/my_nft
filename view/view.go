package view

import (
	"fmt"
	"text/template"
)

var Templates = template.Must(template.ParseFiles("templates/login/adminLogin.html",
	"templates/login/consumerLogin.html", "templates/login/sellerLogin.html"))

func Page(fname string) *template.Template {
	tf, err := template.ParseFiles("templates/"+fname+".html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(err)
		tf, _ = template.New("index").Parse("<html><body><h1>No Template</h1></body></html>")
	}
	return tf
}
