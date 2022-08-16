package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nft-site/lib"
	"path/filepath"
)

type image struct {
	Code string `json:"code"`
}

var images []image

type response struct {
	Images []image `json:"images"`
	Width  int     `json:"width"`
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
		img = image{Code: code}
		images = append(images, img)
	}

	// 画像のwidthを取得
	imgPath := filepath.Join("uploaded", fn)
	width, err := lib.GetWidth(imgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := response{Images: images, Width: width}

	// JSONに変換
	b, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	fmt.Printf("b: %s\n", b)
	images = nil
}
