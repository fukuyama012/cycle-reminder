package services

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"math/rand"
)

// RandString ランダムに指定文字数分の文字列を生成する
func RandString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetDB DB接続取得
func GetDB() *gorm.DB {
	return models.DB
}

// InitDB DB初期化
func InitDB()  {
	if GetDB() == nil {
		models.InitDB()
	}
}