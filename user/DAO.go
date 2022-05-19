package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	*MySQLConf `yaml:"MySQL"`
	*RedisConf `yaml:"Redis"`
}

type MySQLConf struct {
	User   string `yaml:"User"`
	Passwd string `yaml:"Passwd"`
	Net    string `yaml:"Net"`
	Addr   string `yaml:"Addr"`
	DBName string `yaml:"DBName"`
}

type RedisConf struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

var config *Config
var db *sql.DB
var rdb *redis.Client
var ctx context.Context

func init() {
	yamlFile, err := ioutil.ReadFile("user/config/config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	initMySQL()
	initRedis()
}

func initMySQL() {
	cfg := mysql.Config{
		User:   config.MySQLConf.User,
		Passwd: config.MySQLConf.Passwd,
		Net:    config.MySQLConf.Net,
		Addr:   config.MySQLConf.Addr,
		DBName: config.MySQLConf.DBName,
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
		Addr:     config.RedisConf.Addr,
		Password: config.RedisConf.Password, // no password set
		DB:       config.RedisConf.DB,       // use default DB
	})
	ctx = context.Background()
	val, err := rdb.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val)
	}
}
