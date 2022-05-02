package user

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "lx2001812xx",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "wxproject_dev",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(nil)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func newUser(userID string, userName string) error {
	_, err := db.Exec("INSERT INTO user (userID, userName) VALUE (?, ?)", userID, userName)
	if err != nil {
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func selectUser(userID string) (User, error) {
	var res User
	row := db.QueryRow("SELECT * FROM user WHERE userID = ?", userID)
	if err := row.Scan(&res.userID, &res.userName); err != nil {
		if err == sql.ErrNoRows {
			return res, fmt.Errorf("albumsById %s: no such album", userID)
		}
		return res, fmt.Errorf("albumsById %s: %v", userID, err)
	}
	return res, nil
}
