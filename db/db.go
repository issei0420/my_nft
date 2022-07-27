package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

func RegisterDb(r *http.Request) ([]string, error) {

	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("RegisterDB: %v", err)
	}
	fv := map[string]string{} // Form Values
	for k, v := range r.PostForm {
		fv[k] = v[0]
	}

	pbyte := []byte(fv["password"])
	pHash := sha512.Sum512(pbyte)
	xpHash := fmt.Sprintf("%x", pHash)

	if fv["userType"] == "出品者" {

		evls, err := isExist("sellers", []string{"nickname", "mail"}, []string{fv["nickname"], fv["mail"]})
		if err != nil {
			return evls, fmt.Errorf("RegisterDB: %v", err)
		}
		if len(evls) >= 1 {
			return evls, nil
		}

		_, err = db.Exec("INSERT INTO sellers (family_name, first_name, nickname, company, mail, password) VALUES(?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["mail"], xpHash)
		if err != nil {
			return evls, fmt.Errorf("RegisterDB: %v", err)
		}
		return evls, nil

	} else {

		evls, err := isExist("consumers", []string{"nickname", "mail"}, []string{fv["nickname"], fv["mail"]})
		if err != nil {
			return evls, fmt.Errorf("RegisterDB: %v", err)
		}
		if len(evls) >= 1 {
			fmt.Println("len(evls)")
			return evls, nil
		}

		_, err = db.Exec("INSERT INTO consumers (family_name, first_name, nickname, company, lottery_units, mail, password) VALUES(?, ?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["lotteryUnits"], fv["mail"], xpHash)
		if err != nil {
			return evls, fmt.Errorf("RegisterDB: %v", err)
		}
		return evls, nil
	}
}

func isExist(t string, cls []string, vls []string) ([]string, error) {
	var d int
	var evls []string // Exist Values
	for i, c := range cls {
		sq := fmt.Sprintf("SELECT id from %s where %s = ?", t, c)
		row := db.QueryRow(sq, vls[i])
		if err := row.Scan(&d); err != nil {
			if err == sql.ErrNoRows {

			} else {
				return evls, fmt.Errorf("isExist: %v", err)
			}
		} else {
			evls = append(evls, vls[i])
		}
	}
	return evls, nil
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
	Filename   string
	UploadDate string
}

func GetAllImages() ([]Image, error) {
	var imgs []Image
	rows, err := db.Query("SELECT file_name, created_at FROM images")
	if err != nil {
		return imgs, fmt.Errorf("GetAllImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Filename, &img.UploadDate); err != nil {
			return nil, fmt.Errorf("getAllImages: %v", err)
		}
		imgs = append(imgs, img)
	}
	if err := rows.Err(); err != nil {
		return imgs, fmt.Errorf("getAllImages: %v", err)
	}
	return imgs, nil
}
