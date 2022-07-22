package view

import (
	"fmt"
	"text/template"
)

// view
// type Temps struct {
// 	notemp *template.Template
// 	indx   *template.Template
// 	login  *template.Template
// }

// ユーザごとに違うヘッダーを作る
// ユーザーの種類ごとにmustParse

func Page(fname string) *template.Template {
	tf, err := template.ParseFiles("templates/"+fname+".html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(err)
		tf, _ = template.New("index").Parse("<html><body><h1>No Template</h1></body></html>")
	}
	return tf
}
