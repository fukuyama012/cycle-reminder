package models

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB()  {
	c := mysql.Config{
		DBName:               os.Getenv("MYSQL_DATABASE"),
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Addr:                 os.Getenv("MYSQL_ADDRESS")+":"+os.Getenv("MYSQL_PORT"),
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := gorm.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Panicf("Failed gorm.Open %v\n", err)
	}

	db.DB()
	db.AutoMigrate(User{})
	DB = db
}

func CloseDB()  {
	if err := DB.Close(); err != nil {
		log.Panicf("Failed gorm.Close %v\n", err)
	}
}