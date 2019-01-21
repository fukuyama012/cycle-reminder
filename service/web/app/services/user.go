package services

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
)

// GetUserIDOrCreateUserID 登録Eメールが有ればUserIdを取得、無ければ登録してUserIDを取得
func GetUserIDOrCreateUserID(email string) (uint, error) {
	user := models.User{}
	err := user.GetByEmail(models.DB, email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// DBエラー
			return 0, err
		}
		// ユーザー未登録
		return createUser(email)
	}
	return user.ID, nil
}

// CheckUserID ユーザー登録チェック
func CheckUserID(id uint) (*models.User, error) {
	user := models.User{}
	err := user.GetById(models.DB, id)
	if err != nil {
		// record not found or err
		return nil, err
	}
	return &user, nil
}

// ユーザー登録
func createUser(email string) (uint, error) {
	user, err := models.CreateUser(models.DB, email)
	if err != nil {
		// 失敗
		return 0, err
	}
	return user.ID, nil
}
