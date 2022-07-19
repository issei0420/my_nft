package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nft-site/controller"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// database handle
	var db *sql.DB
	connectDB(db)

	dir, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.HandleFunc("/", controller.Index())
	http.HandleFunc("/login", controller.Login("consumer"))
	http.HandleFunc("/admin/login", controller.Login("admin"))
	http.HandleFunc("/seller/login", controller.Login("seller"))
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

func connectDB(db *sql.DB) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "nft",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
