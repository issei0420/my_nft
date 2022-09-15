package handler

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"nft-site/db"
	"nft-site/lib"
	"nft-site/view"
	"os"

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
		ses.Values["login"], ses.Values["mail"], ses.Values["userType"] = nil, nil, nil

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
		ses.Values["login"], ses.Values["mail"], ses.Values["userType"] = nil, nil, nil
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
				return
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
			ses.Values["mail"] = string(accnt)
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
	err := view.AdminParse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Item struct {
		Consumers []db.Consumer
		Sellers   []db.Seller
		UserType  string
	}

	cons, err := db.GetAllConsumers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slrs, err := db.GetAllSellers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item := Item{Consumers: cons, Sellers: slrs, UserType: utype}

	err = view.AdminTemps.ExecuteTemplate(w, "usrList.html", item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager(w, r, "admin")

	type Item struct {
		UserType string
		ResMap   map[string]int
		Form     url.Values
	}

	if err := view.AdminParse(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// ユニーク制限のある項目の値に重複がないかチェック
		UKMap := map[string]string{"mail": r.Form["mail"][0], "nickname": r.Form["nickname"][0]}
		resMap, err := db.UniqueCheck(r.Form["table"][0], UKMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重複があればリダイレクト
		for _, v := range resMap {
			if v == 0 {
				item := Item{
					UserType: "admin",
					ResMap:   resMap,
					Form:     r.Form,
				}
				err := view.AdminTemps.ExecuteTemplate(w, "register.html", item)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		// Form Values
		fv := map[string]string{}
		for k, v := range r.Form {
			fv[k] = v[0]
		}
		// パスワードのハッシュ化
		xpHash := lib.MakeHash(fv["password"])
		// 会員情報の登録
		if err := db.RegisterDb(fv, xpHash); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/usrList", http.StatusFound)
	} else {
		item := Item{UserType: "admin", Form: nil}
		err := view.AdminTemps.ExecuteTemplate(w, "register.html", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var u db.User
		id := r.Form["id"][0]
		table := r.Form["table"][0]
		if r.Form["yaunTable"][0] == "consumers" {
			var c db.Consumer
			c, err := db.GetConsumerFromId(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u = db.User{
				Consumer: c,
			}
		} else {
			var s db.Seller
			s, err := db.GetSellerFromId(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			u = db.User{
				Seller: s,
			}
		}

		// 重複チェック
		diffMap := lib.CheckDiff(u, r.Form)
		UKMap := make(map[string]string)
		if 1 <= len(diffMap) {
			for k := range diffMap {
				if k == "nickname" || k == "mail" {
					UKMap[k] = diffMap[k]
				}
			}
		}

		if 1 <= len(UKMap) {
			resMap, err := db.UniqueCheck(table, UKMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			type Item struct {
				Id       string
				UserType string
				ResMap   map[string]int
				Form     url.Values
			}
			// 重複があればリダイレクト
			for _, v := range resMap {
				if v == 0 {
					item := Item{
						Id:       "id",
						UserType: "admin",
						ResMap:   resMap,
						Form:     r.Form,
					}
					err := view.AdminTemps.ExecuteTemplate(w, "edit.html", item)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}
		if 1 <= len(diffMap) {
			err := db.UpdateUser(diffMap, id, table)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/usrList", http.StatusFound)

	} else {
		sessionManager(w, r, "admin")
		utype := r.FormValue("utype")
		id := r.FormValue("id")

		var u db.User
		if utype == "consumer" {
			var c db.Consumer
			c, err := db.GetConsumerFromId(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u.Consumer = c
		} else {
			var s db.Seller
			s, err := db.GetSellerFromId(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u.Seller = s
		}

		var ic bool
		if u.Consumer.Id != 0 {
			ic = true
		} else {
			ic = false
		}

		item := struct {
			User       db.User
			UserType   string
			Form       string
			IsConsumer bool
		}{
			User:       u,
			UserType:   "admin",
			Form:       "",
			IsConsumer: ic,
		}

		err := view.AdminTemps.ExecuteTemplate(w, "edit.html", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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

		// copy file
		path := fmt.Sprintf("uploaded/%s", fn)
		fd, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fd.Close()
		io.Copy(fd, f)

		// split image
		err = lib.SplitImage(fn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
	mail, utype := sessionManager(w, r, "seller")
	var imgs []db.Image
	var err error
	if utype == "seller" {
		imgs, err = db.GetSellerImages(mail)
	} else if utype == "admin" {
		imgs, err = db.GetAllImages()
	} else {
		log.Fatal("ImgListHandler: invalid usertype")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item := struct {
		UserType string
		Images   []db.Image
	}{
		UserType: utype,
		Images:   imgs,
	}
	err = view.AdminTemps.ExecuteTemplate(w, "imgList.html", item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	_, utype := sessionManager(w, r, "seller")
	fn := r.FormValue("filename")
	item := struct {
		UserType string
		FileName string
	}{
		UserType: utype,
		FileName: fn,
	}

	err := view.AdminTemps.ExecuteTemplate(w, "image.html", item)
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
			return "", ""
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
