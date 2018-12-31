package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setUpBefore()
	ret := m.Run()
	tearDownAfter()
	os.Exit(ret)
}

func setUpBefore()  {
	models.InitDB()
}

func tearDownAfter()  {
}
