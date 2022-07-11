package main

import (
	"net/http"
	"nft-site/controller"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.HandleFunc("/", controller.Index())
	http.HandleFunc("/login", controller.Login())
	http.HandleFunc("/register", controller.Register())
	http.HandleFunc("/upload", controller.Upload())
	http.HandleFunc("/imgList", controller.ImgList())
	http.HandleFunc("/image", controller.Image())
	http.HandleFunc("/lottery", controller.Lottery())
	http.HandleFunc("/result", controller.Result())
	http.HandleFunc("/myImgList", controller.MyImgList())
	http.HandleFunc("/myImage", controller.MyImage())
	http.ListenAndServe("", nil)
}
