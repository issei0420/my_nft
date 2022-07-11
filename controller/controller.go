package controller

import (
	"log"
	"net/http"
	"nft-site/view"
)

func Index() func(w http.ResponseWriter, rq *http.Request) {
	// switch controller depending on user's type
	hh := UsrList()
	return hh
}

func Login() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("login")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func UsrList() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("usrList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func Register() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("register")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func Upload() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("upload")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func ImgList() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("imgList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func Image() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("image")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func Lottery() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("lottery")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func Result() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("result")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func MyImgList() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("myImgList")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}

func MyImage() func(w http.ResponseWriter, rq *http.Request) {
	tf := view.Page("myImage")
	hh := func(w http.ResponseWriter, rq *http.Request) {
		er := tf.Execute(w, nil)
		if er != nil {
			log.Fatal(er)
		}
	}
	return hh
}
