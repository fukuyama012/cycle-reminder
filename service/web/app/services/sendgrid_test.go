package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Eメール送信
func TestSendMail(t *testing.T) {
	// CIでテストし辛いのでAPI KEY未設定時は実メール送信はスキップしておく
	if len(os.Getenv("SENDGRID_API_KEY")) == 0 {
		return
	}
	
	tests := []struct {
		ToEmail string
		Subject string
		Content string
	}{
		{"test@example.com", "subject", "text contents"},
		{"test@example.com", "サブジェクト", "コンテンツです"},
	}
	for _, tt := range tests {
		response, err :=  services.SendMail(tt.ToEmail, tt.Subject, tt.Content)
		assert.NoError(t, err)
		assert.True(t, services.IsSuccessStatusCode(response))
	}
}

// Eメール送信　バリデーションエラー
func TestSendMailValidationError(t *testing.T) {
	tests := []struct {
		ToEmail string
		Subject string
		Content string
	}{
		{"", "subject", "text コンテンツです"}, // Eメール空
		{"@gmail.com", "subject", "text コンテンツです"}, // Eメール形式不正
		{"test@example.com", "", "コンテンツ"}, // メールタイトル無し
		{"test@example.com", "サブジェクト", ""}, // 本文コンテンツ無し
		{"test@example.com", services.RandString(101), ""}, // メールタイトル文字数超え
		{"test@example.com", "sub", services.RandString(1001)}, // 本文コンテンツ文字数超え
	}
	for _, tt := range tests {
		response, err :=  services.SendMail(tt.ToEmail, tt.Subject, tt.Content)
		assert.Error(t, err)
		assert.Nil(t, response)
	}
}

// レスポンスが成功ステータスか？
func TestIsSuccessStatusCode(t *testing.T) {
	tests := []struct {
		StatusCode int
		Return bool
	}{
		{200, true},
		{202, true},
		{400, false},
		{500, false},
	}
	for _, tt := range tests {
		response := &rest.Response{
			StatusCode: tt.StatusCode,
		}
		assert.Equal(t, tt.Return, services.IsSuccessStatusCode(response))
	}
}