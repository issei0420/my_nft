package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nft-site/db"
	"nft-site/lib"
	"path/filepath"
)

type image struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
}

var images []image

type response struct {
	Images   []image      `json:"images"`
	Width    int          `json:"width"`
	Portions []db.Portion `json:"portions"`
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	fn := r.FormValue("filename")
	name := filepath.Base(fn[:len(fn)-len(filepath.Ext(fn))])
	var part, path, code string
	var img image
	var err error
	for i := 0; i < 100; i++ {
		part = fmt.Sprintf("No_%d.png", i)
		path = filepath.Join("out/original", name, part)
		code, err = lib.Encode(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		img = image{Id: i, Code: code}
		images = append(images, img)
	}

	// 画像のwidthを取得
	imgPath := filepath.Join("uploaded", fn)
	width, err := lib.GetWidth(imgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 抽選済み部の情報を取得
	portions, err := db.GetSoldPortions(fn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := response{
		Images:   images,
		Width:    width,
		Portions: portions,
	}

	// JSONに変換
	b, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	images = nil
}

type data struct {
	Password string
	Table    string
	Id       string
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var d data
	json.NewDecoder(r.Body).Decode(&d)

	// パスワードのハッシュ化
	xpHash := lib.MakeHash(d.Password)
	// パスワードの更新
	var message string
	err := db.UpdatePassword(xpHash, d.Table, d.Id)
	if err != nil {
		message = "パスワードの更新に失敗しました"
	} else {
		message = "パスワードを変更しました"
	}
	b, err := json.MarshalIndent(message, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func GetPortionTotal(w http.ResponseWriter, r *http.Request) {
	var id []string
	json.NewDecoder(r.Body).Decode(&id)
	res, err := db.GetPortionTotal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
