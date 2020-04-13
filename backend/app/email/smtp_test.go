package email

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_SMTPSend(t *testing.T) {
	fromEmail := os.Getenv("SMTP_FROM")
	toEmail := os.Getenv("SMTP_TO")
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	useTls, _ := strconv.ParseBool(os.Getenv("SMTP_TLS"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if fromEmail == "" || toEmail == "" || host == "" || port == 0 {
		t.Skip("SMTP_FROM, SMTP_TO, SMTP_HOST or SMTP_PORT is empty, skip the SMTPSend test ...")
	}

	sndr := NewSMTPSender(&SmtpParams{
		Host:     host,
		Port:     port,
		TLS:      useTls,
		Username: username,
		Password: password,
		TimeOut:  3 * time.Second,
	}, nil)

	subject := "Sending with SMTP is Not safe"

	sndr.SetSubject(subject)
	sndr.SetFrom(fromEmail)

	htmlContent := "<strong>and may cause source IP leaking problem</strong>"
	t.Logf("try send via SMTP from %s to %s", fromEmail, toEmail)

	msg, err := sndr.(*SMTPSender).BuildMessage(toEmail, htmlContent, "text/html")
	if err != nil {
		t.Error(err)
	}
	t.Logf("mail msg: %s", msg)
	if !strings.Contains(msg, htmlContent) {
		t.Errorf("BuildMessage lost body")
	}

	// after the BuildMessage test, append the extra line now,
	// because if we include it in htmlContent it will formatted to
	//<strong>and may cause source IP leaking problem</strong><p>this is a testin=
	//        g email from CI</p>
	htmlContent += "<p>this is a testing email from CI</p>"
	err = sndr.Send(toEmail, htmlContent)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("mail send sucess")
	}
}
