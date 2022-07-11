package view

import (
	"text/template"
)

// view
// type Temps struct {
// 	notemp *template.Template
// 	indx   *template.Template
// 	login  *template.Template
// }

// create template file
func UsrList() *template.Template {
	tf, er := template.ParseFiles("static/templates/usrList.html", "static/templates/header.html")
	if er != nil {
		tf, _ = template.New("index").Parse("<html><body><h1>No Template</h1></body></html>")
	}
	return tf
}
