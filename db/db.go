package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

func RegisterDb(r *http.Request) (int64, []string, error) {
	// columns which already exist
	cols := []string{}

	if err := r.ParseForm(); err != nil {
		return 0, cols, fmt.Errorf("RegisterDB: %v", err)
	}
	fv := map[string]string{}
	for k, v := range r.PostForm {
		fv[k] = v[0]
	}

	if fv["userType"] == "出品者" {

		res1, err := isExist("sellers", "nickname", fv["nickname"])
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		if res1 {
			cols = append(cols, fv["nickname"])
		}

		res2, err := isExist("sellers", "mail", fv["mail"])
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		if res2 {
			cols = append(cols, fv["mail"])
		}

		// return if either column exists
		if res1 || res2 {
			return 0, cols, nil
		}

		result, err := db.Exec("INSERT INTO sellers (family_name, first_name, nickname, company, mail, password) VALUES(?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["mail"], fv["password"])
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		return id, cols, nil

	} else {

		res1, err := isExist("consumers", "nickname", fv["nickname"])
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		if res1 {
			cols = append(cols, fv["nickname"])
		}

		res2, err := isExist("consumers", "mail", fv["mail"])
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		if res2 {
			cols = append(cols, fv["mail"])
		}

		// return if either column exists
		if res1 || res2 {
			return 0, cols, nil
		}
		result, err := db.Exec("INSERT INTO consumers (family_name, first_name, nickname, company, lottery_units, mail, password) VALUES(?, ?, ?, ?, ?, ?, ?)",
			fv["familyName"], fv["firstName"], fv["nickname"], fv["company"], fv["lotteryUnits"], fv["mail"], fv["password"])
		if err != nil {
			return 0, cols, fmt.Errorf("eRegisterDB: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, cols, fmt.Errorf("RegisterDB: %v", err)
		}
		return id, cols, nil
	}
}

func isExist(t string, c string, v string) (bool, error) {
	var i int
	sq := fmt.Sprintf("SELECT id from %s where %s = ?", t, c)
	row := db.QueryRow(sq, v)
	if err := row.Scan(&i); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return true, err
		}
	}
	return true, nil
}
