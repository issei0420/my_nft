package controller

import (
	"log"
	"net/http"
	"nft-site/view"
)

// controller
func UsrListController() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.UsrList()
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func IndexSwitcher() func(w http.ResponseWriter, rq *http.Request) {
	// switch controller depending on user's type
	hh := UsrListController()
	return hh
}
