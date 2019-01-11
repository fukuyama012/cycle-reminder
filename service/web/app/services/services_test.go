package services_test

import (
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

func tearDownAfter()  {
}
