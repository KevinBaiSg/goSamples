package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("filepath directory error ", err)
		return
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("viper ReadInConfig error ", err)
		return
	}

	url 		:= viper.GetString("url")
	dbname 		:= viper.GetString("dbname")
	user 		:= viper.GetString("user")
	password 	:= viper.GetString("password")

	dsn := strings.Join([]string{user, ":", password,
		"@tcp(", url, ")/", dbname, "?charset=utf8"}, "")
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil{
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connnect success")
}
