package handler

import (
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	_, utype := sessionManager(w, r, "seller")
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
		Images   []db.ImageNames
		Data     db.Data
		Array    [100]int // 抽選口数0-100
	}

	var imgs []db.ImageNames
	imgs, err := db.GetAllImageNames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var array [100]int
	for i := 0; i < 100; i++ {
		array[i] = i + 1
	}

	if err := view.AdminParse(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var d db.Data
		json.NewDecoder(r.Body).Decode(&d)

		// ユニーク制限のある項目の値に重複がないかチェック
		UKMap := map[string]string{"mail": d.Mail, "nickname": d.Nickname}

		resMap, err := db.UniqueCheck(d.UserType, UKMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重複があればフォームを再表示
		var invalid bool
		for _, v := range resMap {
			if v == 0 {
				invalid = true
				break
			}
		}

		type response struct {
			Invalid int            `json:"invalid"`
			ResMap  map[string]int `json:"resmap"`
		}

		if invalid {
			res := response{
				Invalid: 1,
				ResMap:  resMap,
			}
			b, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		// パスワードのハッシュ化
		xpHash := lib.MakeHash(d.Password)
		// テーブルへの追加
		id, err := db.RegisterDb(d.FamilyName, d.FirstName, d.Nickname, d.Mail, d.Company, d.UserType, xpHash)
		if err != nil {
			log.Fatal((err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if d.UserType == "sellers" {
			http.Redirect(w, r, "/usrList", http.StatusFound)
		}
		// ticketsテーブルへの追加
		if err := db.InsertTickets(id, d.ImageUnits); err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := response{
			Invalid: 0,
		}
		b, err := json.MarshalIndent(res, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)

	} else {
		item := Item{
			UserType: "admin",
			Images:   imgs,
			Array:    array,
			ResMap:   map[string]int{"mail": 1, "nickname": 1},
		}
		err = view.AdminTemps.ExecuteTemplate(w, "register.html", item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {

	var imgs []db.ImageNames
	imgs, err := db.GetAllImageNames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var array [100]int
	for i := 0; i < 100; i++ {
		array[i] = i + 1
	}

	if r.Method == "POST" {
		var d db.Data
		json.NewDecoder(r.Body).Decode(&d)
		fmt.Printf("d: %v\n", d)
		var u db.User
		id := d.Id
		table := d.UserType
		UKMap := make(map[string]string)

		if table == "consumers" {
			var c db.Consumer
			c, err := db.GetConsumerFromId(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u = db.User{
				Consumer: c,
			}
			// ニックネームやメールアドレスが変更されている場合はUKMapに追加し、あとで重複を確認
			if u.Consumer.Nickname != d.Nickname {
				UKMap["nickname"] = d.Nickname
			}
			if u.Consumer.Mail != d.Mail {
				UKMap["mail"] = d.Mail
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
			// ニックネームやメールアドレスが変更されている場合はUKMapに追加し、あとで重複を確認
			if u.Seller.Nickname != d.Nickname {
				UKMap["nickname"] = d.Nickname
			}
			if u.Seller.Mail != d.Mail {
				UKMap["mail"] = d.Mail
			}
		}

		// ユニーク制限のある項目の値に重複がないかチェック
		resMap, err := db.UniqueCheck(table, UKMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 重複があればフォームにエラーを表示
		var invalid bool
		for _, v := range resMap {
			if v == 0 {
				invalid = true
				break
			}
		}

		type response struct {
			Invalid int            `json:"invalid"`
			ResMap  map[string]int `json:"resmap"`
		}

		if invalid {
			res := response{
				Invalid: 1,
				ResMap:  resMap,
			}
			b, err := json.MarshalIndent(res, "", "\t")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		res := response{
			Invalid: 0,
		}
		b, err := json.MarshalIndent(res, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)

		err = db.UpdateUser(d, id, table)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			fmt.Printf("c(handler): %v\n", c)
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
			Images     []db.ImageNames
			Array      [100]int
		}{
			User:       u,
			UserType:   "admin",
			Form:       "",
			IsConsumer: ic,
			Images:     imgs,
			Array:      array,
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
		if err := view.SellerParse(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	if utype == "admin" {
		err = view.AdminTemps.ExecuteTemplate(w, "imgList.html", item)
	} else if utype == "seller" {
		err = view.SellerTemps.ExecuteTemplate(w, "imgList.html", item)
	} else {
		log.Fatal()
	}

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

		var images []db.LotteryImage
		images, err := db.GetImages(mail)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		item := struct {
			UserType string
			Images   []db.LotteryImage
		}{
			UserType: utype,
			Images:   images,
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
