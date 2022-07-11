package main

import (
	"net/http"
	"nft-site/controller"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.HandleFunc("/index", controller.IndexSwitcher())
	http.ListenAndServe("", nil)
}
