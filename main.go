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
	http.HandleFunc("/login", handler.LoginHandler)            // seller & consumer
	http.HandleFunc("/admin/login", handler.AdminLoginHandler) // admin
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

// func write() {
// 	pbyte := []byte("Th5RLynP")
// 	pHash := sha512.Sum512(pbyte)
// 	xpHash := fmt.Sprintf("%x", pHash)

// 	err := os.WriteFile("data", []byte(xpHash), 0600)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("File written!")
// }
