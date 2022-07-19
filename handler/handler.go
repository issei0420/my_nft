package handler

import (
	"fmt"
	"log"
	"net/http"
	"nft-site/db"
	"nft-site/view"
	"text/template"
)

func IndexHandler() func(w http.ResponseWriter, rq *http.Request) {
	// switch controller depending on user's type
	hh := UsrListHandler()
	return hh
}

// consumer
func LoginHandler(utype string) func(w http.ResponseWriter, rq *http.Request) {
	var tf *template.Template
	switch utype {
	case "consumer":
		tf = view.Page("login/consumerLogin")
	case "admin":
		tf = view.Page("login/adminLogin")
	case "seller":
		tf = view.Page("login/sellerLogin")
	}
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func UsrListHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("usrList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func RegisterHandler() func(w http.ResponseWriter, r *http.Request) {
	var tf *template.Template
	hh := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id, cols, err := db.RegisterDb(r)
			if err != nil {
				log.Fatal(err)
			}
			if len(cols) > 1 {
				fmt.Println(cols)
				tf = view.Page("register")
				er := tf.Execute(w, nil)
				if er != nil {
					log.Fatal(er)
				}
				fmt.Println("cols が存在する")
			} else {
				tf = view.Page("usrList")
				er := tf.Execute(w, nil)
				if er != nil {
					log.Fatal(er)
				}
				fmt.Println("colsがない")
			}
			fmt.Printf("Successully inserted! id: %d", id)
		} else {
			tf = view.Page("usrList")
			er := tf.Execute(w, nil)
			if er != nil {
				log.Fatal(er)
			}
		}
	}
	return hh
}

func UploadHander() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("upload")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func ImgListHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("imgList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func ImageHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("image")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func LotteryHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("lottery")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func ResultHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("result")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func MyImgListHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("myImgList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func MyImageHandler() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("myImage")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}
