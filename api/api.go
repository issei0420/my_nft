package api

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type image struct {
	Code string `json:"code"`
}

var images []image

func GetImages(c *gin.Context) {
	data := encode()
	img := image{Code: data}
	images = append(images, img)
	c.IndentedJSON(http.StatusOK, images)
}

//エンコード
func encode() string {
	path := filepath.Join("uploaded/rena.png")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat() //FileInfo interface
	if err != nil {
		log.Fatal(err)
	}
	size := fi.Size() //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data)
}

func StartApi() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/images", GetImages)

	router.Run("localhost:8080")
}
