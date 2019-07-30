package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"gopkg.in/testfixtures.v2"
	"log"
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
}

func prepareTestDB() {
	fixtures, err := testfixtures.NewFolder(models.GetDB().DB(), &testfixtures.MySQL{}, "../../tests/fixtures/models")
	if err != nil {
		log.Fatal(err)
	}
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func tearDownAfter()  {
	models.CloseDB()
}
