package main

import (
	"net/http"
	"nft-site/db"
	"nft-site/handler"
	"os"
)

func main() {
	// database handle
	db.ConnectDb()

	dir, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.HandleFunc("/", handler.IndexHandler())
	http.HandleFunc("/login", handler.LoginHandler("consumer"))
	http.HandleFunc("/admin/login", handler.LoginHandler("admin"))
	http.HandleFunc("/seller/login", handler.LoginHandler("seller"))
	http.HandleFunc("/register", handler.RegisterHandler())
	http.HandleFunc("/upload", handler.UploadHander())
	http.HandleFunc("/imgList", handler.ImgListHandler())
	http.HandleFunc("/image", handler.ImageHandler())
	http.HandleFunc("/lottery", handler.LotteryHandler())
	http.HandleFunc("/result", handler.ResultHandler())
	http.HandleFunc("/myImgList", handler.MyImgListHandler())
	http.HandleFunc("/myImage", handler.MyImageHandler())
	http.ListenAndServe(":8080", nil)
}
