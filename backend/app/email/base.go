package email

type BaseSender struct {
	From        string
	Subject     string
	Headers     map[string]string
	ContentType string // text/plain or text/html
}

var _ IBaseSender = &BaseSender{}

func (s *BaseSender) AddHeader(header, value string) {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}
	s.Headers[header] = value
}

func (s *BaseSender) ResetHeaders() {
	s.Headers = nil
}

func (s *BaseSender) SetFrom(from string) {
	s.From = from
}

func (s *BaseSender) SetSubject(subject string) {
	s.Subject = subject
}
