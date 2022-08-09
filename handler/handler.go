package handler

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"nft-site/db"
	"nft-site/lib"
	"nft-site/view"
	"os"
	"text/template"

	"github.com/gorilla/sessions"
)

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte("nft-session-key"))

type Message struct {
	Msg string
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// get session
	ses, err := cs.Get(r, "login-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check user type
	utype := ses.Values["userType"].(string)
	switch utype {
	case "consumer":
		http.Redirect(w, r, "/lottery", http.StatusFound)
	case "seller":
		http.Redirect(w, r, "/upload", http.StatusFound)
	case "admin":
		http.Redirect(w, r, "/usrList", http.StatusFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	_, utype := sessionManager(w, r, "admin")

	view.AdminParse()

	type Item struct {
		Consumers []db.Consumers
		Sellers   []db.Sellers
		UserType  string
	}

	cons, err := db.GetAllConsumers()
	if err != nil {
		log.Fatal(err)
	}
	slrs, err := db.GetAllSellers()
	if err != nil {
		log.Fatal(err)
	}
	item := Item{Consumers: cons, Sellers: slrs, UserType: utype}

	err = view.AdminTemps.ExecuteTemplate(w, "usrList.html", item)
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
	if r.Method == "POST" {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get file
		f, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// check if its unique
		fn := r.FormValue("file-name")
		exist, err := db.IsExistImage(fn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exist {
			item := struct {
				UserType string
				Message  string
			}{
				UserType: "seller",
				Message:  "同名の画像がすでに出品されています。",
			}
			view.SellerTemps.ExecuteTemplate(w, "upload.html", item)
			return
		}

		// split image
		err = lib.SplitImage(fn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// insert into images
		id, err := db.InsertImage(fn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// insert into upload
		ses, _ := cs.Get(r, "login-session")
		mail := ses.Values["mail"].(string)
		err = db.InsertUpload(mail, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// copy file
		path := fmt.Sprintf("upload/%s", fn)
		fd, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fd.Close()
		io.Copy(fd, f)

		http.Redirect(w, r, "/imgList", http.StatusFound)

	} else {
		_, utype := sessionManager(w, r, "seller")
		view.SellerParse()
		item := struct {
			UserType string
			Message  string
		}{
			UserType: utype,
			Message:  "",
		}
		err := view.SellerTemps.ExecuteTemplate(w, "upload.html", item)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ImgListHandler(w http.ResponseWriter, r *http.Request) {
	_, utype := sessionManager(w, r, "seller")
	imgs, err := db.GetAllImages()
	if err != nil {
		log.Fatal(err)
	}
	item := struct {
		UserType string
		Images   []db.Image
	}{
		UserType: utype,
		Images:   imgs,
	}
	err = view.SellerTemps.ExecuteTemplate(w, "imgList.html", item)
	if err != nil {
		log.Fatal(err)
	}
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "seller")
	err := view.SellerTemps.ExecuteTemplate(w, "image.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LotteryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		imageId, units := r.FormValue("imageId"), r.FormValue("units")
		ses, _ := cs.Get(r, "login-session")
		mail := ses.Values["mail"].(string)
		// get random portion
		soldP, err := db.GetPortion(imageId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		randP, err := lib.RandomPortion(soldP, units)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// insert into Lottery
		id, err := db.InsertLottery(imageId, mail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// insert into portion
		err = db.InsertPortion(id, randP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// update lottery_units
		db.UpdateUnits(mail, len(randP))
		item := struct {
			Units    int
			UserType string
			ImageId  string
		}{
			Units:    len(randP),
			UserType: "consumer",
			ImageId:  imageId,
		}
		err = view.ConsumerTemps.ExecuteTemplate(w, "result.html", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		mail, utype := sessionManager(w, r, "consumer")
		view.ConsumerParse()
		// get all images
		imgs, err := db.GetAllImages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// get Units
		units, err := db.GetUnits(mail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var unitsSlice []int
		for i := 0; i < units; i++ {
			unitsSlice = append(unitsSlice, i+1)
		}
		item := struct {
			UserType string
			Images   []db.Image
			Units    []int
		}{
			UserType: utype,
			Images:   imgs,
			Units:    unitsSlice,
		}
		err = view.ConsumerTemps.ExecuteTemplate(w, "lottery.html", item)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func MyImgListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	} else {
		mail, utype := sessionManager(w, r, "consumer")
		myImgs, err := db.GetAllMyImages(mail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item := struct {
			UserType string
			MyImgs   []db.MyImage
		}{
			UserType: utype,
			MyImgs:   myImgs,
		}
		err = view.ConsumerTemps.ExecuteTemplate(w, "myImgList.html", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func MyImageHandler(w http.ResponseWriter, r *http.Request) {
	mail, _ := sessionManager(w, r, "consumer")
	id := r.FormValue("id")
	fn := r.FormValue("filename")
	if fn == "" {
		var err error
		fn, err = db.GetFilenameFromId(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	getPs, err := db.GetMyPortion(id, mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	getP := make(map[uint8]struct{})
	for _, v := range getPs {
		getP[uint8(v)] = struct{}{}
	}

	err = lib.ProcessImage(fn, getP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	item := struct {
		UserType string
		FileName string
	}{
		UserType: "consumer",
		FileName: fn,
	}
	err = view.ConsumerTemps.ExecuteTemplate(w, "myImage.html", item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sessionManager(w http.ResponseWriter, r *http.Request, u string) (string, string) {
	// get session
	ses, err := cs.Get(r, "login-session")
	if err != nil {
		log.Fatal(err)
	}
	// check login status
	flg := ses.Values["login"].(bool)
	if !flg {
		switch utype := ses.Values["userType"].(string); utype {
		case "admin":
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		case "seller":
			http.Redirect(w, r, "/login?usr=seller", http.StatusFound)
		case "consumer":
			http.Redirect(w, r, "/login", http.StatusFound)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
	// check user type
	utype := ses.Values["userType"].(string)
	if utype != u && utype != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	// get id
	mail := ses.Values["mail"].(string)
	return mail, utype
}
