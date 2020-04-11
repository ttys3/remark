package email

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MailgunSender(t *testing.T) {
	fromEmail := os.Getenv("MAILGUN_FROM")
	toEmail := os.Getenv("MAILGUN_TO")
	dm := os.Getenv("MAILGUN_DOMAIN")
	key := os.Getenv("MAILGUN_API_KEY")
	if fromEmail == "" || toEmail == "" || dm == "" || key == "" {
		t.Skip("MAILGUN_FROM, MAILGUN_TO, MAILGUN_DOMAIN or MAILGUN_API_KEY is empty, skip the MailgunSender test ...")
	}
	sndr := NewMailgunSender(dm, key, 0)
	assert.Equal(t, DefaultEmailTimeout, sndr.(*MailgunSender).Timeout, "test default timeout value")

	subject := "Sending with Mailgun is Fun"
	sndr.SetSubject(subject)
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong><p>this is a testing email from CI</p>"

	t.Logf("try send from %s to %s", fromEmail, toEmail)
	err := sndr.Send(toEmail, htmlContent)
	if err != nil {
		t.Error(err)
	}
}