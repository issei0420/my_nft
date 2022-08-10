package main

import (
	"log"
	"net/http"
	"nft-site/api"
	"nft-site/db"
	"nft-site/handler"
	"os"
)

func main() {
	db.ConnectDb()
	api.StartApi()

	dir, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.Handle("/out/", http.StripPrefix("/out/", http.FileServer(http.Dir(dir+"/out/"))))
	http.Handle("/uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir(dir+"/uploaded/"))))

	http.HandleFunc("/login", handler.LoginHandler)            // seller & consumer
	http.HandleFunc("/admin/login", handler.AdminLoginHandler) // admin
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/usrList", handler.UsrListHandler)
	http.HandleFunc("/register", handler.Register())
	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/imgList", handler.ImgListHandler)
	http.HandleFunc("/image", handler.ImageHandler)
	http.HandleFunc("/lottery", handler.LotteryHandler)
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/myImgList", handler.MyImgListHandler)
	http.HandleFunc("/myImage", handler.MyImageHandler)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal()
	}
}
