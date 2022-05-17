package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB
var rdb *redis.Client

func init() {
	initMySQL()
	initRedis()
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
		Addr:     "175.24.163.131:6380",
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
