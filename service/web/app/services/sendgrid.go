package services

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

const (
	SENDGRID_RESPONSE_STATUS_CODE_OK = 200
	SENDGRID_RESPONSE_STATUS_CODE_ACCEPTED = 202
)

type Message struct {
	ToEmail string `validate:"required,email"`
	Subject string `validate:"max=100"`
	Content string `validate:"required,max=1000"`
}

func (mes Message) validate() error {
	return validator.New().Struct(mes)
}

func SendMail(toEmail, subject, content string) (*rest.Response, error) {
	message, err := makeMessage(toEmail, subject, content)
	if err != nil {
		return nil, err
	}
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	return client.Send(message)
}

// 成功ステータスコードか？
func IsSuccessStatusCode(response *rest.Response) bool {
	switch response.StatusCode {
	case SENDGRID_RESPONSE_STATUS_CODE_OK:
	case SENDGRID_RESPONSE_STATUS_CODE_ACCEPTED:
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
	// plain, html同じ内容（暫定）
	message := mail.NewSingleEmail(from, subject, to, content, content)
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
