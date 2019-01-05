package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_GetById(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
		out string
	}{
		{1, "test1@example.com"},
	}
	for _, tt := range tests {
		user := models.User{}
		if err := user.GetById(tt.in); err != nil{
			t.Error(err)
		}
		if user.Email != tt.out {
			assert.Equal(t, tt.out, user.Email)
		}
	}
	}

func TestUser_GetByIdNotFound(t *testing.T)  {
	prepareTestDB()
	var assertT = assert.New(t)
	tests := []struct {
		in  uint
	}{
		{999},
		{12345},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetById(tt.in); 
		assertT.Equal(gorm.ErrRecordNotFound, err)
	}
}

func TestUser_GetAll(t *testing.T) {
	prepareTestDB()
	users, err := models.GetAllUsers()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(users))
}
