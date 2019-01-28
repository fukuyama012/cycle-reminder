package services

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

const (
	sendgridResponseStatusCodeOk       = 200
	sendgridResponseStatusCodeAccepted = 202
)

// Message is struct
type Message struct {
	ToEmail string `validate:"required,email"`
	Subject string `validate:"required,max=100"`
	Content string `validate:"required,max=1000"`
}

func (mes Message) validate() error {
	return validator.New().Struct(mes)
}

// SendMail 送信内容バリデーションしてEメール送信
func SendMail(toEmail, subject, content string) (*rest.Response, error) {
	message, err := makeMessage(toEmail, subject, content)
	if err != nil {
		return nil, err
	}
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	return client.Send(message)
}

// IsSuccessStatusCode 成功ステータスコードか？
func IsSuccessStatusCode(response *rest.Response) bool {
	switch response.StatusCode {
	case sendgridResponseStatusCodeOk, sendgridResponseStatusCodeAccepted:
		return true
	}
	return false
}

func makeMessage(toEmail, subject, content string) (*mail.SGMailV3, error) {
	if err := validateMessage(toEmail, subject, content); err != nil {
		return nil, err
	}

	from := mail.NewEmail("", os.Getenv("SENDGRID_FROM_MAIL"))
	to := mail.NewEmail("", toEmail)
	// 改行が消えるのでtext/htmlは設定しない
	plainText := mail.NewContent("text/plain", content)
	message := mail.NewV3MailInit(from, subject, to, plainText)
	return message, nil
}

func validateMessage(toEmail, subject, content string) error {
	mes := Message{
		ToEmail: toEmail,
		Subject: subject,
		Content: content,
	}
	// validator.v9
	if err := mes.validate(); err != nil {
		return err
	}
	return nil
}

