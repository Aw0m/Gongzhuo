package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"log"
	"wxProjectDev/utils"
)

var db *sql.DB
var rdb *redis.Client

func init() {
	initMySQL()
	initRedis()
}

func createUser(userID string, userName string) error {
	_, err := db.Exec("INSERT INTO user (userID, userName) VALUE (?, ?)", userID, userName)
	if err != nil {
		log.Println("create user时，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func selectUser(userID string) (User, error) {
	var res User
	row := db.QueryRow("SELECT * FROM user WHERE userID = ?", userID)
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
	_, err := db.Exec("UPDATE user SET userName = ? WHERE userID = ?", userName, userID)
	if err != nil {
		log.Println("update user时，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func createTeam(userID, teamName string) error {
	utils.SetMachineId(0)
	teamID := utils.GetSnowflakeId()
	_, err := db.Exec("INSERT INTO team (teamID, teamName, creatorID) VALUE (?, ?, ?)", teamID, teamName, userID)
	if err != nil {
		log.Println("create team，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func updateTeam(teamID int64, teamName string) error {
	_, err := db.Exec("UPDATE team SET teamName = ? WHERE teamID = ?", teamName, teamID)
	if err != nil {
		log.Println("update team，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func initMySQL() {
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

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "175.24.163.131:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	ctx := context.Background()
	val, err := rdb.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val)
	}
}
