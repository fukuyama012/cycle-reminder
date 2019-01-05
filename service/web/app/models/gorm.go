package models

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB()  {
	connectDB()
	migrate()
}

func connectDB()  {
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
	DB = db
}

func migrate()  {
	// User
	DB.AutoMigrate(User{})
	// ReminderSetting
	DB.AutoMigrate(ReminderSetting{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	DB.Model(&ReminderSetting{}).AddUniqueIndex("reminder_settings_unq_user_id_number", "user_id", "number")
	// ReminderSchedule
	DB.AutoMigrate(ReminderSchedule{}).AddForeignKey("reminder_setting_id", "reminder_settings(id)", "CASCADE", "RESTRICT")
	DB.Model(&ReminderSchedule{}).AddUniqueIndex("reminder_schedules_unq_reminder_setting_id", "reminder_setting_id")
	// ReminderLog 
	// Userが消えたらCASCADE削除
	DB.AutoMigrate(ReminderLog{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT") 
	// Settingを消してもCASCADE削除しない
	DB.Model(&ReminderLog{}).AddForeignKey("reminder_setting_id", "reminder_settings(id)", "RESTRICT", "RESTRICT") 
}

func CloseDB()  {
	if err := DB.Close(); err != nil {
		log.Panicf("Failed gorm.Close %v\n", err)
	}
}