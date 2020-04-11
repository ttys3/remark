// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
package email

import (
	"os"
	"testing"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

func Test_SendgridSender(t *testing.T) {
	fromEmail := os.Getenv("SENDGRID_FROM")
	toEmail := os.Getenv("SENDGRID_TO")
	key := os.Getenv("SENDGRID_API_KEY")

	if fromEmail == "" || toEmail == "" || key == "" {
		t.Skip("SENDGRID_FROM, SENDGRID_TO or SENDGRID_API_KEY is empty, skip the SendgridSender test ...")
	}

	sndr := NewSendgridSender(key, 0)
	assert.Equal(t, DefaultEmailTimeout, sndr.(*SendgridSender).TimeOut, "test default timeout value")

	subject := "Sending with SendGrid is Fun"
	sndr.SetSubject(subject)
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong><p>this is a testing email from CI</p>"

	t.Logf("try send from %s to %s", fromEmail, toEmail)
	err := sndr.Send(toEmail, htmlContent)
	if err != nil {
		t.Error(err)
	}
}

func Test_SendgridSDK(t *testing.T) {
	fromEmail := os.Getenv("SENDGRID_FROM")
	toEmail := os.Getenv("SENDGRID_TO")
	key := os.Getenv("SENDGRID_API_KEY")

	if fromEmail == "" || toEmail == "" || key == "" {
		t.Skip("SENDGRID_FROM, SENDGRID_TO or SENDGRID_API_KEY is empty, skip the SendgridSDK test ...")
	}

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	from := mail.NewEmail("Example User", fromEmail)
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", toEmail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong><p>this is a testing email from CI</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	t.Logf("try send from %s to %s", fromEmail, toEmail)
	response, err := client.Send(message)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("StatusCode: %#v, Body: %#v, Headers: %#v", response.StatusCode, response.Body, response.Headers)
	}
}