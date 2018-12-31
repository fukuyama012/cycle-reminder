package models_test

import (
	"database/sql"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/testfixtures.v2"
	"log"
	"os"
	"testing"
)

var (
	db_fixture *sql.DB
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
	prepareDB()
	prepareFixtures()
}

func prepareDB() {
	c := mysql.Config{
		DBName:               os.Getenv("MYSQL_DATABASE"),
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Addr:                 os.Getenv("MYSQL_ADDRESS")+":"+os.Getenv("MYSQL_PORT"),
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db_, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	db_fixture = db_
}

func prepareFixtures()  {
	fixtures_, err := testfixtures.NewFolder(db_fixture, &testfixtures.MySQL{}, "../../tests/fixtures")
	if err != nil {
		log.Fatal(err)
	}
	testfixtures.SkipDatabaseNameCheck(true)
	fixtures = fixtures_
}

func prepareTestDB() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func tearDownAfter()  {
	if err := db_fixture.Close(); err != nil {
		log.Fatal(err)
	}
	models.CloseDB()
}
