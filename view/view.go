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

// ユーザごとに違うヘッダーを作る

func Page(fname string) *template.Template {
	tf, er := template.ParseFiles("static/templates/"+fname+".html",
		"static/templates/header.html", "static/templates/footer.html")
	if er != nil {
		tf, _ = template.New("index").Parse("<html><body><h1>No Template</h1></body></html>")
	}
	return tf
}
