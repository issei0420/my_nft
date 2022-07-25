package view

import (
	"fmt"
	"text/template"
)

var LoginTemps = template.Must(template.ParseFiles("templates/login/adminLogin.html",
	"templates/login/consumerLogin.html", "templates/login/sellerLogin.html"))

var AdminTemps *template.Template

func AdminParse() error {
	tmps, err := template.ParseFiles("templates/header.html", "templates/footer.html",
		"templates/image.html", "templates/imgList.html", "templates/register.html",
		"templates/upload.html", "templates/usrList.html")
	if err != nil {
		return fmt.Errorf("AdminRender: %v", err)
	}
	AdminTemps = tmps
	return nil
}

func Page(fname string) *template.Template {
	tf, err := template.ParseFiles("templates/"+fname+".html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(err)
		tf, _ = template.New("index").Parse("<html><body><h1>No Template</h1></body></html>")
	}
	return tf
}
