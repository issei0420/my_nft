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

	"github.com/gorilla/sessions"
)

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte("nft-session-key"))

type Message struct {
	Msg string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("table")
	if r.Method == "POST" {
		// reset session
		ses, err := cs.Get(r, "login-session")
		if err != nil {
			log.Fatal(err)
		}
		ses.Values["login"] = nil

		p, err := db.GetPassword(t, r.FormValue("mail"))
		var m Message
		if err != nil {
			if err == sql.ErrNoRows {
				m.Msg = "アカウントが存在しません"
			} else {
				log.Fatal(err)
			}
			if t == "sellers" {
				view.LoginTemps.ExecuteTemplate(w, "sellerLogin.html", m)
			} else {
				view.LoginTemps.ExecuteTemplate(w, "consumerLogin.html", m)
			}
		} else {
			pbyte := []byte(r.FormValue("password"))
			pHash := sha512.Sum512(pbyte)
			xpHash := fmt.Sprintf("%x", pHash)
			if xpHash != p {
				m.Msg = "パスワードが間違っています"
				if t == "sellers" {
					view.LoginTemps.ExecuteTemplate(w, "sellerLogin.html", m)
				} else {
					view.LoginTemps.ExecuteTemplate(w, "consumerLogin.html", m)
				}
			} else {
				// 認証成功
				ses.Values["login"] = true
				ses.Values["mail"] = r.FormValue("mail")
				if t == "sellers" {
					ses.Values["userType"] = "seller"
					ses.Save(r, w)
					http.Redirect(w, r, "/upload", http.StatusFound)
				} else {
					ses.Values["userType"] = "consumer"
					ses.Save(r, w)
					http.Redirect(w, r, "/lottery", http.StatusFound)
				}
			}
		}
	} else {
		usr := r.FormValue("usr")
		if usr == "seller" {
			err := view.LoginTemps.ExecuteTemplate(w, "sellerLogin.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := view.LoginTemps.ExecuteTemplate(w, "consumerLogin.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// reset session
		ses, err := cs.Get(r, "login-session")
		if err != nil {
			log.Fatal(err)
		}
		ses.Values["login"] = nil
		// check ID
		accnt, err := os.ReadFile("data/accnt.txt")
		if err != nil {
			log.Fatal(err)
		}

		var m Message

		if r.FormValue("mail") != string(accnt) {
			m.Msg = "IDが間違っています"
			err := view.LoginTemps.ExecuteTemplate(w, "adminLogin.html", m)
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
			err = view.LoginTemps.ExecuteTemplate(w, "adminLogin.html", m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			ses.Values["login"] = true
			ses.Values["userType"] = "admin"
			if ses.Save(r, w) != nil {
				log.Fatal(ses.Save(r, w))
			}
			http.Redirect(w, r, "/usrList", http.StatusFound)
		}
	} else {
		err := view.LoginTemps.ExecuteTemplate(w, "adminLogin.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ses, err := cs.Get(r, "login-session")
	if err != nil {
		log.Fatal()
	}
	ses.Values["login"] = false
	err = ses.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}
	utype := ses.Values["userType"].(string)
	switch utype {
	case "admin":
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	case "seller":
		http.Redirect(w, r, "/login?usr=seller", http.StatusFound)
	default:
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func UsrListHandler(w http.ResponseWriter, r *http.Request) {
	// check login status
	ses, err := cs.Get(r, "login-session")
	if err != nil {
		log.Fatal(err)
	}
	flg := ses.Values["login"].(bool)
	if !flg {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	view.AdminParse()

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

	err = view.AdminTemps.ExecuteTemplate(w, "usrList.html", usrs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "seller")
	view.SellerParse()
	err := view.ConsumerTemps.ExecuteTemplate(w, "upload.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ImgListHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "seller")
	err := view.ConsumerTemps.ExecuteTemplate(w, "imgList.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "seller")
	err := view.ConsumerTemps.ExecuteTemplate(w, "image.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func LotteryHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "consumer")
	view.ConsumerParse()
	err := view.ConsumerTemps.ExecuteTemplate(w, "lottery.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "consuemer")
	err := view.ConsumerTemps.ExecuteTemplate(w, "result.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func MyImgListHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "consumer")
	err := view.ConsumerTemps.ExecuteTemplate(w, "myImgList.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func MyImageHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "consumer")
	err := view.ConsumerTemps.ExecuteTemplate(w, "myImage.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func sessionManager(w http.ResponseWriter, r *http.Request, u string) string {
	// get session
	ses, err := cs.Get(r, "login-session")
	if err != nil {
		log.Fatal(err)
	}
	// check login status
	flg := ses.Values["login"].(bool)
	if !flg {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	// check user type
	utype := ses.Values["userType"].(string)
	if utype != u && utype != "admin" {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// get id
	mail := ses.Values["mail"].(string)
	return mail
}
