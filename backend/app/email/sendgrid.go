// SendGrid(https://sendgrid.com) Trial Plan provides 40,000 emails for 30 days
// After your trial ends, you can send 100 emails/day for free

package email

import (
	"fmt"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MailgunConfig contain settings for mailgun API
type SendgridSender struct {
	sg      *sendgrid.Client
	APIKey  string        // the SendGrid API key
	Timeout time.Duration // TCP connection timeout
	BaseSender
}

func NewSendgridSender(apiKey string, timeout time.Duration) EmailSender {
	if timeout == 0 {
		timeout = DefaultEmailTimeout
	}
	sender := &SendgridSender{
		APIKey:  apiKey,
		Timeout: timeout,
	}

	// Create an instance of the sendgrid Client
	sender.sg = sendgrid.NewSendClient(apiKey)
	return sender
}

func (s *SendgridSender) Name() string {
	return "sendgrid"
}

func (s *SendgridSender) Send(to, text string) error {
	if s.From == "" {
		return fmt.Errorf("sendgrid: empty From. the from object must be provided for every email send")
	}
	if to == "" {
		return fmt.Errorf("sendgrid: empty to. at least one receipt should be provided")
	}
	fromEmail := mail.NewEmail("", s.From)
	toEmail := mail.NewEmail("", to)
	sgmail := mail.NewSingleEmail(fromEmail, s.Subject, toEmail, text, text)

	// extra headers used mainly for List-Unsubscribe feature
	// see more info via https://sendgrid.com/docs/ui/sending-email/list-unsubscribe/
	if s.Headers != nil && len(s.Headers) > 0 {
		sgmail.Headers = s.Headers
	}
	s.SetTimeout(s.Timeout)
	resp, err := s.sg.Send(sgmail)
	if err != nil {
		return fmt.Errorf("sendgrid: request failed: %w", err)
	}
	// 2xx responses indicate a successful request
	// see https://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html
	if resp.StatusCode%100 != 2 {
		return fmt.Errorf("sendgrid: send failed with err: %+v", resp.Body)
	}
	log.Infof("sendgrid: send to %s success, subject: %s, StatusCode: %d\n", to, s.Subject, resp.StatusCode)
	return nil
}

func (s *SendgridSender) SetTimeout(timeout time.Duration) {
	if timeout != 0 {
		s.Timeout = timeout
	}
	sendgrid.DefaultClient.HTTPClient.Timeout = s.Timeout
}

// String representation of Email object
func (s *SendgridSender) String() string {
	return fmt.Sprintf("email.sender.sendgrid: API %s", "v3")
}
