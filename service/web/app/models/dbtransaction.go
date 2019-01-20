package models

import "github.com/jinzhu/gorm"

// トランザクション実行し情報を受取る
func TransactAndReceiveData(db *gorm.DB, txFunc func(*gorm.DB) (interface{}, error)) (data interface{}, err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	data, err = txFunc(tx)
	return
}

// トランザクション実行
func Transact(db *gorm.DB, txFunc func(*gorm.DB) (error)) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	err = txFunc(tx)
	return
}