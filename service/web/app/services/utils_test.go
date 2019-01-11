package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUtils_RandString(t *testing.T) {
	tests := []struct {
		in  int
		out int
	}{
		{1, 1},
		{10, 10},
		{999, 999},
	}
	for _, tt := range tests {
		str := services.RandString(tt.in)
		// 桁数チェック
		assert.Equal(t, tt.out, len(str))
		// 全て英字である事
		assert.True(t, check_regexp("^[a-zA-Z]+$", str))
	}
}

func check_regexp(reg, str string) bool {
	return (regexp.MustCompile(reg).Match([]byte(str)))
}
