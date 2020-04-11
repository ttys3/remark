// Mailgun(https://www.mailgun.com/) Free Plan provides 10,000 Emails per month

package email

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunConfig contain settings for mailgun API
type MailgunSender struct {
	mg          *mailgun.MailgunImpl
	Domain      string
	APIKey      string
	Timeout     time.Duration // TCP connection timeout
	From        string
	Subject     string
	Headers     map[string]string
	ContentType string // text/plain or text/html
}

func NewMailgunSender(domain, apiKey string, timeout time.Duration) EmailSender {
	if timeout == 0 {
		timeout = DefaultEmailTimeout
	}
	sender := &MailgunSender {
		Domain:  domain,
		APIKey:  apiKey,
		Timeout: timeout,
	}

	// Create an instance of the Mailgun Client
	sender.mg = mailgun.NewMailgun(domain, apiKey)
	return sender
}

func (s *MailgunSender) Name() string {
	return "mailgun"
}

func (s *MailgunSender) Send(to, text string) error {
	if s.From == "" {
		return fmt.Errorf("mailgun: empty From. the from object must be provided for every email send")
	}
	if to == "" {
		return fmt.Errorf("mailgun: empty to. at least one receipt should be provided")
	}
	message := s.mg.NewMessage(s.From, s.Subject, text, to)
	message.SetHtml(text)

	// extra headers used mainly for List-Unsubscribe feature
	// You can enable Mailgun’s Unsubscribe functionality by turning it on in the settings area for your domain.
	// Mailgun can automatically provide an unsubscribe footer in each email you send.
	// Mailgun will automatically prevent future emails being sent to recipients that have unsubscribed.
	// You can edit the unsubscribed address list from your Control Panel or through the API.
	// see more info via https://documentation.mailgun.com/en/latest/api-unsubscribes.html
	// and https://documentation.mailgun.com/en/latest/user_manual.html#tracking-unsubscribes
	if s.Headers != nil && len(s.Headers) > 0{
		keys := make([]string, 0, len(s.Headers))
		for k := range s.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			message.AddHeader(k, s.Headers[k])
		}
	}
	s.SetTimeout(s.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultEmailTimeout)
	defer cancel()
	// Send the message	with a 10 second timeout
	resp, id, err := s.mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("mailgun: send failed: %w", err)
	}
	fmt.Printf("mailgun: send to %s success, ID: %s Resp: %s\n", to, id, resp)
	return nil
}

func (s *MailgunSender) AddHeader(header, value string) {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}
	s.Headers[header] = value
}

func (s *MailgunSender) ResetHeaders() {
	s.Headers = nil
}

func (s *MailgunSender) SetFrom(from string) {
	s.From = from
}

func (s *MailgunSender) SetSubject(subject string) {
	s.Subject = subject
}

func (s *MailgunSender) SetTimeout(timeout time.Duration) {
	if timeout != 0 {
		s.Timeout = timeout
	}
	s.mg.Client().Timeout = s.Timeout
}

// String representation of Email object
func (s *MailgunSender) String() string {
	return fmt.Sprintf("email.sender.mailgrun: domain %s", s.Domain)
}