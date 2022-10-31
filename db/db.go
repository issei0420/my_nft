package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

// database handle
var db *sql.DB

func ConnectDb() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "nft",
		AllowNativePasswords: true,
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

func RegisterDb(familyName, firstName, nickname, company, mail, userType, xpHash string) error {
	// Sellers
	if userType == "sellers" {
		_, err := db.Exec("INSERT INTO sellers (family_name, first_name, nickname, company, mail, password) VALUES(?, ?, ?, ?, ?, ?)",
			familyName, firstName, nickname, company, mail, xpHash)
		if err != nil {
			return fmt.Errorf("RegisterDB: %v", err)
		}
		return nil
		// Consumers
	} else if userType == "consumers" {
		_, err := db.Exec("INSERT INTO consumers (family_name, first_name, nickname, company, mail, password) VALUES(?, ?, ?, ?, ?, ?)",
			familyName, firstName, nickname, company, mail, xpHash)
		if err != nil {
			return fmt.Errorf("RegisterDB: %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("RegisterDb: no such table %s", userType)
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

type Consumer struct {
	Id           int
	FamilyName   string
	FirstName    string
	Nickname     string
	Mail         string
	Company      string
	LotteryUnits string
	ImgNum       string
}

type Seller struct {
	Id         int
	FamilyName string
	FirstName  string
	Nickname   string
	Mail       string
	Company    string
	UploadNum  string
}

type User struct {
	Consumer
	Seller
}

func GetConsumerFromId(id string) (Consumer, error) {
	var c Consumer
	row := db.QueryRow("SELECT id, family_name, first_name, nickname, mail, company, lottery_units FROM consumers where id = ?", id)
	err := row.Scan(&c.Id, &c.FamilyName, &c.FirstName, &c.Nickname, &c.Mail, &c.Company, &c.LotteryUnits)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("GetConsumerFromId %s: No such consumer", id)
		}
		return c, fmt.Errorf("GetConsumerFromId %s: %v", id, err)
	}
	return c, nil
}

func GetAllConsumers() ([]Consumer, error) {
	var cons []Consumer
	rows, err := db.Query("select c.id, c.family_name, c.first_name, c.nickname, c.mail, c.company, count(portion) " +
		"from lottery l right join consumers c on l.consumer_id = c.id " +
		"left join portion p on l.id = p.lottery_id group by c.id")
	if err != nil {
		return cons, fmt.Errorf("getAllConsumers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var con Consumer
		if err := rows.Scan(&con.Id, &con.FamilyName, &con.FirstName, &con.Nickname, &con.Mail, &con.Company, &con.LotteryUnits, &con.ImgNum); err != nil {
			return nil, fmt.Errorf("getAllConsumers: %v", err)
		}
		cons = append(cons, con)
	}
	if err := rows.Err(); err != nil {
		return cons, fmt.Errorf("getAllConsumers: %v", err)
	}
	return cons, nil
}

func UpdateUser(diffMap map[string]string, id, table string) error {
	st := ""
	for k, v := range diffMap {
		st += fmt.Sprintf(" %s = \"%s\",", k, v)
	}
	st = st[:len(st)-1]
	sq := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, st)
	_, err := db.Exec(sq, id)
	if err != nil {
		return fmt.Errorf("UpdateConsumer: %v", err)
	}
	return nil
}

func GetAllSellers() ([]Seller, error) {
	var slrs []Seller
	rows, err := db.Query("select s.id, s.family_name, s.first_name, s.nickname, s.mail, s.company, count(u.id) from sellers s left join upload u on s.id = u.seller_id group by s.id;")
	if err != nil {
		return slrs, fmt.Errorf("getAllSellers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var slr Seller
		if err := rows.Scan(&slr.Id, &slr.FamilyName, &slr.FirstName, &slr.Nickname, &slr.Mail, &slr.Company, &slr.UploadNum); err != nil {
			return nil, fmt.Errorf("getAllSellers: %v", err)
		}
		slrs = append(slrs, slr)
	}
	if err := rows.Err(); err != nil {
		return slrs, fmt.Errorf("getAllSellers: %v", err)
	}
	return slrs, nil
}

func GetSellerFromId(id string) (Seller, error) {
	var s Seller
	row := db.QueryRow("SELECT id, family_name, first_name, nickname, mail, company FROM sellers where id = ?", id)
	err := row.Scan(&s.Id, &s.FamilyName, &s.FirstName, &s.Nickname, &s.Mail, &s.Company)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, fmt.Errorf("GetSellerFromId %s: No such seller", id)
		}
		return s, fmt.Errorf("GetSellerFromId %s: %v", id, err)
	}
	return s, nil
}

func GetPassword(table, mail string) (string, error) {
	var p string
	sq := fmt.Sprintf("SELECT password from %s where mail = ?", table)
	row := db.QueryRow(sq, mail)
	if err := row.Scan(&p); err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", fmt.Errorf("GetPassword: %v", err)
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
	Id           string
	SellerFamily string
	SellerFirst  string
	Filename     string
	Count        int
	UploadDate   string
}

func GetSellerImages(mail string) ([]Image, error) {
	var imgs []Image
	rows, err := db.Query("SELECT i.file_name, COUNT(portion), i.created_at FROM lottery l "+
		"INNER JOIN portion p ON l.id = p.lottery_id "+
		"RIGHT JOIN images i ON l.image_id = i.id "+
		"LEFT JOIN upload u ON i.id = u.image_id "+
		"INNER JOIN sellers s ON u.seller_id = s.id "+
		"GROUP BY i.file_name, s.mail "+
		"HAVING s.mail=?", mail)
	if err != nil {
		return imgs, fmt.Errorf("GetSellerImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Filename, &img.Count, &img.UploadDate); err != nil {
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
	rows, err := db.Query("SELECT A.id, A.file_name, s.family_name, s.first_name, A.count, u.created_at FROM " +
		"( SELECT i.id, i.file_name, COUNT(portion) AS count, i.created_at FROM lottery l " +
		"RIGHT JOIN images i ON l.image_id = i.id LEFT JOIN portion p ON l.id = p.lottery_id " +
		"LEFT JOIN upload u ON i.id = u.image_id LEFT JOIN sellers s ON u.seller_id = s.id GROUP BY i.id) AS A " +
		"LEFT JOIN upload u ON A.id = u.image_id LEFT JOIN sellers s ON u.seller_id = s.id")
	if err != nil {
		return imgs, fmt.Errorf("GetAllImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Id, &img.Filename, &img.SellerFamily, &img.SellerFirst, &img.Count, &img.UploadDate); err != nil {
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
	Id           int
	Filename     string
	CountPortion string
}

func GetAllMyImages(mail string) ([]MyImage, error) {
	var myImgs []MyImage
	rows, err := db.Query("SELECT i.id, i.file_name, count(portion) FROM lottery l INNER JOIN  images i ON l.image_id = i.id INNER JOIN consumers c ON l.consumer_id = c.id inner join portion p on l.id = p.lottery_id WHERE c.mail = ? group by i.id", mail)
	if err != nil {
		return nil, fmt.Errorf("GetAllMyImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var myImg MyImage
		if err := rows.Scan(&myImg.Id, &myImg.Filename, &myImg.CountPortion); err != nil {
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

func UpdatePassword(p, t, id string) error {
	sq := fmt.Sprintf("UPDATE %s SET password = ? WHERE id = ? ", t)
	_, err := db.Exec(sq, p, id)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %v", err)
	}
	return nil
}

type Total struct {
	Id       int
	FileName string
	Count    int
}

func GetPortionTotal() ([]Total, error) {
	var everyTotal []Total

	rows, err := db.Query("select l.consumer_id, i.file_name, count(portion)  from lottery l inner join consumers c on l.consumer_id = c.id inner join images i on l.image_id = i.id inner join portion p on l.id = p.lottery_id group by l.consumer_id, l.image_id")
	if err != nil {
		return nil, fmt.Errorf("GetPortionTotal: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t Total
		if err := rows.Scan(&t.Id, &t.FileName, &t.Count); err != nil {
			return nil, fmt.Errorf("GetPortionTotal: %v", err)
		}
		everyTotal = append(everyTotal, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPortionTotal: %v", err)
	}
	return everyTotal, nil
}

type Uploaded struct {
	Id       int
	FileName string
}

func GetUploadedImages() ([]Uploaded, error) {
	var allUploaded []Uploaded
	rows, err := db.Query("select u.seller_id, i.file_name from upload u right join sellers s on u.seller_id = s.id inner join images i on u.image_id = i.id;")
	if err != nil {
		return nil, fmt.Errorf("GetUploadedImages: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u Uploaded
		if err := rows.Scan(&u.Id, &u.FileName); err != nil {
			return nil, fmt.Errorf("GetUploadedImages: %v", err)
		}
		allUploaded = append(allUploaded, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUploadeImages: %v", err)
	}
	return allUploaded, nil
}

type ImageNames struct {
	Id       int
	FileName string
}

func GetAllImageNames() ([]ImageNames, error) {
	var imgs []ImageNames
	rows, err := db.Query("SELECT id, file_name FROM images")
	if err != nil {
		return nil, fmt.Errorf("GetAllImageNames: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var img ImageNames
		if err := rows.Scan(&img.Id, &img.FileName); err != nil {
			return nil, fmt.Errorf("GetAllImageNames: %v", err)
		}
		imgs = append(imgs, img)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllImageNames: %v", err)
	}
	return imgs, nil
}
