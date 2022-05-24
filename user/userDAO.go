package user

import (
	"database/sql"
	"fmt"
	"log"
	"wxProjectDev/public"
)

func createUser(userID string, userName string) error {
	_, err := public.DB.Exec("INSERT INTO user (userID, userName) VALUE (?, ?)", userID, userName)
	if err != nil {
		log.Println("create user时，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func selectUser(userID string) (User, error) {
	var res User
	row := public.DB.QueryRow("SELECT * FROM user WHERE userID = ?", userID)
	if err := row.Scan(&res.userID, &res.userName); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no such userID: %s\n", userID)
			return res, err
		}
		return res, fmt.Errorf("selectUser %s: %v", userID, err)
	}
	return res, nil
}

func updateUser(userID string, userName string) error {
	_, err := public.DB.Exec("UPDATE user SET userName = ? WHERE userID = ?", userName, userID)
	if err != nil {
		log.Println("update user时，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}
