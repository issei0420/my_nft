package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"crypto/sha512"

	"github.com/go-sql-driver/mysql"
)

// database handle
var db *sql.DB

func ConnectDb() {
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
	fmt.Println("Database Connected!")
}

func RegisterDb(form url.Values) error {
	// Form Values
	fv := map[string]string{}
	for k, v := range form {
		fv[k] = v[0]
	}
	// encoded password
	pbyte := []byte(fv["password"])
	pHash := sha512.Sum512(pbyte)
	xpHash := fmt.Sprintf("%x", pHash)
	// Sellers
	if fv["table"] == "sellers" {
		_, err := db.Exec("INSERT INTO sellers (family_name, first_name, nickname, company, mail, password) VALUES(?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["mail"], xpHash)
		if err != nil {
			return fmt.Errorf("RegisterDB: %v", err)
		}
		return nil
		// Consumers
	} else if fv["table"] == "consumers" {
		_, err := db.Exec("INSERT INTO consumers (family_name, first_name, nickname, company, lottery_units, mail, password) VALUES(?, ?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["lotteryUnits"], fv["mail"], xpHash)
		if err != nil {
			return fmt.Errorf("RegisterDB: %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("RegisterDb: no such table %s", fv["table"])
	}
}

func UniqueCheck(table string, UKMap map[string]string) (map[string]int, error) {
	var id int
	resMap := make(map[string]int)
	for k, v := range UKMap {
		query := fmt.Sprintf("SELECT id FROM %s WHERE %s = ?", table, k)
		row := db.QueryRow(query, v)
		if err := row.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				resMap[k] = 1 // is unique
				continue
			}
			return resMap, fmt.Errorf("UniqueCheck: %v", err)
		}
		resMap[k] = 0 // not unique
	}
	return resMap, nil
}

type Consumers struct {
	FamilyName string
	FirstName  string
	Nickname   string
	Mail       string
	Company    string
	// haveNum    int
}

func GetAllConsumers() ([]Consumers, error) {
	var cons []Consumers
	rows, err := db.Query("SELECT family_name, first_name, nickname, mail, company FROM consumers")
	if err != nil {
		return cons, fmt.Errorf("getAllConsumers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var con Consumers
		if err := rows.Scan(&con.FamilyName, &con.FirstName, &con.Nickname, &con.Mail, &con.Company); err != nil {
			return nil, fmt.Errorf("getAllConsumers: %v", err)
		}
		cons = append(cons, con)
	}
	if err := rows.Err(); err != nil {
		return cons, fmt.Errorf("getAllConsumers: %v", err)
	}
	return cons, nil
}

type Sellers struct {
	FamilyName string
	FirstName  string
	Nickname   string
	Mail       string
	Company    string
}

func GetAllSellers() ([]Sellers, error) {
	var slrs []Sellers
	rows, err := db.Query("SELECT family_name, first_name, nickname, mail, company FROM sellers")
	if err != nil {
		return slrs, fmt.Errorf("getAllSellers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var slr Sellers
		if err := rows.Scan(&slr.FamilyName, &slr.FirstName, &slr.Nickname, &slr.Mail, &slr.Company); err != nil {
			return nil, fmt.Errorf("getAllSellers: %v", err)
		}
		slrs = append(slrs, slr)
	}
	if err := rows.Err(); err != nil {
		return slrs, fmt.Errorf("getAllSellers: %v", err)
	}
	return slrs, nil
}

func GetPassword(t string, mail string) (string, error) {
	var p string
	sq := fmt.Sprintf("SELECT password from %s where mail = ?", t)
	row := db.QueryRow(sq, mail)
	if err := row.Scan(&p); err != nil {
		return "", err
	}
	return p, nil
}

func InsertImage(fn string) (int64, error) {
	result, err := db.Exec("INSERT INTO images (file_name) VALUES (?)", fn)
	if err != nil {
		return 0, fmt.Errorf("InsertImage: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertImage: %v", err)
	}
	return id, nil
}

func InsertUpload(mail string, id int64) error {
	_, err := db.Exec("Insert into upload (image_id, seller_id) values(?, (select id from sellers where mail = ?))", id, mail)
	if err != nil {
		return fmt.Errorf("InsertUpload: %v", err)
	}
	return nil
}

func GetSellerId(mail string) (int64, error) {
	var id int64
	row := db.QueryRow("SELECT id from sellers where mail = ?", mail)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("GetSellerId %s: No such seller", mail)
		}
		return 0, fmt.Errorf("GetSellerId %s: %v", mail, err)
	}
	return id, nil
}

type Image struct {
	Id         string
	Filename   string
	UploadDate string
}

func GetSellerImages(mail string) ([]Image, error) {
	var imgs []Image
	rows, err := db.Query("SELECT i.id, i.file_name, i.created_at FROM upload u "+
		"INNER JOIN images i ON u.image_id = i.id INNER JOIN sellers s ON u.seller_id = s.id WHERE s.mail = ?", mail)
	if err != nil {
		return imgs, fmt.Errorf("GetSellerImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Id, &img.Filename, &img.UploadDate); err != nil {
			return nil, fmt.Errorf("GetSellerImages: %v", err)
		}
		imgs = append(imgs, img)
	}
	if err := rows.Err(); err != nil {
		return imgs, fmt.Errorf("GetSellerImages: %v", err)
	}
	return imgs, nil
}

func GetAllImages() ([]Image, error) {
	var imgs []Image
	rows, err := db.Query("SELECT id, file_name, created_at FROM images")
	if err != nil {
		return imgs, fmt.Errorf("GetAllImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Id, &img.Filename, &img.UploadDate); err != nil {
			return nil, fmt.Errorf("GetAllImages: %v", err)
		}
		imgs = append(imgs, img)
	}
	if err := rows.Err(); err != nil {
		return imgs, fmt.Errorf("GetAllImages: %v", err)
	}
	return imgs, nil
}

func InsertLottery(imageId string, mail string) (int64, error) {
	result, err := db.Exec("INSERT INTO lottery (image_id, consumer_id)"+
		"VALUES (?, (SELECT id FROM consumers WHERE mail=?))", imageId, mail)
	if err != nil {
		return 0, fmt.Errorf("InsertLottery: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertLottery: %v", err)
	}
	return id, nil
}

func InsertPortion(id int64, portion []uint8) error {
	for _, p := range portion {
		_, err := db.Exec("INSERT INTO portion (lottery_id, portion) VALUES(?, ?)", id, p)
		if err != nil {
			return fmt.Errorf("InsertPortion: %v", err)
		}
	}
	return nil
}

func GetPortion(imageId string) ([]uint8, error) {
	var soldP []uint8
	rows, err := db.Query("SELECT portion FROM lottery INNER JOIN portion ON lottery.id = portion.lottery_id WHERE image_id=?", imageId)
	if err != nil {
		return nil, fmt.Errorf("GetPortion: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p uint8
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("GetPortion: %v", err)
		}
		soldP = append(soldP, p)
	}
	return soldP, nil
}

func GetUnits(mail string) (int, error) {
	row := db.QueryRow("SELECT lottery_units FROM consumers WHERE mail = ?", mail)
	var units int
	if err := row.Scan(&units); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("GetUnits %s: No such consumer", mail)
		}
		return 0, fmt.Errorf("GetUnits %s: %v", mail, err)
	}
	return units, nil
}

func UpdateUnits(mail string, units int) error {
	_, err := db.Exec("UPDATE consumers SET lottery_units = lottery_units - ? WHERE mail=?", units, mail)
	if err != nil {
		return fmt.Errorf("UpdateUnits: %v", err)
	}
	return nil
}

type MyImage struct {
	Id       int
	Filename string
}

func GetAllMyImages(mail string) ([]MyImage, error) {
	var myImgs []MyImage
	rows, err := db.Query("SELECT DISTINCT i.id, i.file_name FROM lottery l INNER JOIN images i ON l.image_id = i.id "+
		"INNER JOIN consumers c ON l.consumer_id = c.id WHERE c.mail= ?", mail)
	if err != nil {
		return nil, fmt.Errorf("GetAllMyImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var myImg MyImage
		if err := rows.Scan(&myImg.Id, &myImg.Filename); err != nil {
			return nil, fmt.Errorf("GetAllMyImages: %v", err)
		}
		myImgs = append(myImgs, myImg)
	}
	if err != nil {
		return nil, fmt.Errorf("GetAllMyImages: %v", err)
	}
	return myImgs, nil
}

func IsExistImage(fn string) (bool, error) {
	row := db.QueryRow("SELECT id FROM images WHERE file_name = ?", fn)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, fmt.Errorf("IsExistImage: %v", err)
	}
	return true, nil
}

func GetMyPortion(id, mail string) ([]uint8, error) {
	var getP []uint8
	rows, err := db.Query("SELECT portion FROM lottery l INNER JOIN portion p ON l.id = p.lottery_id "+
		"INNER JOIN consumers c ON l.consumer_id=c.id WHERE l.image_id = ? and c.mail = ?", id, mail)
	if err != nil {
		return nil, fmt.Errorf("GetMyPortion: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p uint8
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("GetMyPortion: %v", err)
		}
		getP = append(getP, p)
	}
	return getP, nil
}

func GetFilenameFromId(id string) (string, error) {
	row := db.QueryRow("SELECT file_name FROM images WHERE id = ?", id)
	var fn string
	if err := row.Scan(&fn); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("GetFilenameFromId: No such image of id %s", id)
		}
		return "", fmt.Errorf("GetFilenameFromId: %v", err)
	}
	return fn, nil
}

type Portion struct {
	PId        string `json:"pId"`
	UserID     string `json:"userId"`
	FamilyName string `json:"familyName"`
	FirstName  string `json:"firstName"`
	Date       string `json:"date"`
}

func GetSoldPortions(fn string) ([]Portion, error) {
	var portions []Portion
	rows, err := db.Query("select p.portion, c.id, c.family_name, c.first_name, l.updated_at from lottery l "+
		"inner join images i on l.image_id = i.id inner join consumers c on l.consumer_id = c.id "+
		"inner join portion p on l.id = p.lottery_id where i.file_name = ?", fn)
	if err != nil {
		return nil, fmt.Errorf("GetSoldPortions: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Portion
		if err := rows.Scan(&p.PId, &p.UserID, &p.FamilyName, &p.FirstName, &p.Date); err != nil {
			return nil, fmt.Errorf("GetSoldPortions: %v", err)
		}
		portions = append(portions, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSoldPortions: %v", err)
	}
	return portions, nil
}
