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
	http.HandleFunc("/login", handler.LoginHandler)
	// http.HandleFunc("/admin/login", handler.AdminLoginHandler)
	// http.HandleFunc("/seller/login", handler.SellerLoginHandler)
	http.HandleFunc("/usrList", handler.UsrList())
	http.HandleFunc("/register", handler.Register())
	http.HandleFunc("/upload", handler.Upload())
	http.HandleFunc("/imgList", handler.ImgList())
	http.HandleFunc("/image", handler.Image())
	http.HandleFunc("/lottery", handler.Lottery())
	http.HandleFunc("/result", handler.Result())
	http.HandleFunc("/myImgList", handler.MyImgList())
	http.HandleFunc("/myImage", handler.MyImage())
	http.ListenAndServe(":8080", nil)
}
