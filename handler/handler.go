package handler

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nft-site/db"
	"nft-site/view"
	"os"
	"text/template"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("table")
	if r.Method == "POST" {
		p, err := db.GetPassword(t, r.FormValue("mail"))
		var msg string
		if err != nil {
			if err == sql.ErrNoRows {
				msg = "アカウントが存在しません"
			} else {
				log.Fatal(err)
			}
		} else {
			pbyte := []byte(r.FormValue("password"))
			pHash := sha512.Sum512(pbyte)
			xpHash := fmt.Sprintf("%x", pHash)
			if xpHash != p {
				msg = "パスワードが間違っています"
			} else {
				msg = "認証成功"
				if t == "sellers" {
					http.Redirect(w, r, "/upload", http.StatusFound)
				} else {
					http.Redirect(w, r, "/lottery", http.StatusFound)
				}
			}
		}
		fmt.Printf("ConsumerLoginHandler: %v\n", msg)
	} else {
		usr := r.FormValue("usr")
		if usr == "seller" {
			err := view.Templates.ExecuteTemplate(w, "sellerLogin.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := view.Templates.ExecuteTemplate(w, "consumerLogin.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// check ID
		accnt, err := os.ReadFile("data/accnt.txt")
		if err != nil {
			log.Fatal(err)
		}

		type Message struct {
			Msg string
		}
		var m Message

		if r.FormValue("mail") != string(accnt) {
			m.Msg = "IDが間違っています"
			err := view.Templates.ExecuteTemplate(w, "adminLogin.html", m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// check password
		pswd, err := os.ReadFile("data/pswd.txt")
		if err != nil {
			log.Fatal(err)
		}
		pbyte := []byte(r.FormValue("password"))
		pHash := sha512.Sum512(pbyte)
		xpHash := fmt.Sprintf("%x", pHash)
		if xpHash != string(pswd) {
			m.Msg = "パスワードが間違っています"
		}
		err = view.Templates.ExecuteTemplate(w, "adminLogin.html", m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := view.Templates.ExecuteTemplate(w, "adminLogin.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func UsrList() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("usrList")

	hh := func(w http.ResponseWriter, r *http.Request) {

		type Users struct {
			Consumers []db.Consumers
			Sellers   []db.Sellers
		}

		cons, err := db.GetAllConsumers()
		if err != nil {
			log.Fatal(err)
		}
		slrs, err := db.GetAllSellers()
		if err != nil {
			log.Fatal(err)
		}
		usrs := Users{Consumers: cons, Sellers: slrs}

		if err = tf.Execute(w, usrs); err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func Register() func(w http.ResponseWriter, r *http.Request) {
	var tf *template.Template
	hh := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			evls, err := db.RegisterDb(r)
			if err != nil {
				log.Fatal(err)
			}
			if len(evls) > 1 {
				fmt.Println(evls)
				fmt.Println("すでに存在する値です")
				tf = view.Page("register")
				err := tf.Execute(w, nil)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				http.Redirect(w, r, "/usrList", http.StatusFound)
			}
		} else {
			tf = view.Page("register")
			err := tf.Execute(w, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return hh
}

func Upload() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("upload")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func ImgList() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("imgList")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func Image() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("image")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func Lottery() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("lottery")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func Result() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("result")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func MyImgList() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("myImgList")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}

func MyImage() func(w http.ResponseWriter, r *http.Request) {
	tf := view.Page("myImage")
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := tf.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return hh
}
