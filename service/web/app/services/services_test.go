package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"gopkg.in/testfixtures.v2"
	"log"
	"os"
	"testing"
	"time"
)

var (
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {
	setUpBefore()
	ret := m.Run()
	tearDownAfter()
	os.Exit(ret)
}

func setUpBefore()  {
	models.InitDB()
	prepareFixtures()
}

func prepareFixtures()  {
	fixtures_, err := testfixtures.NewFolder(models.DB.DB(), &testfixtures.MySQL{}, "../../tests/fixtures/services")
	if err != nil {
		log.Fatal(err)
	}
	fixtures = fixtures_
}

func prepareTestDB() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
	// CIで投入情報反映がテストに間に合わない(?)場合が有るので少し処理を待つ
	time.Sleep(100 * time.Millisecond)
}

func tearDownAfter()  {
	models.CloseDB()
}
