package email

import (
	"fmt"
	"time"

	"github.com/go-pkgz/auth/provider"
)

type EmailSender interface {
	provider.Sender // implement for github.com/go-pkgz/auth/provider.VerifyHandler
	fmt.Stringer
	SetTimeout(timeout time.Duration)
	Name() string
	IBaseSender
}

type IBaseSender interface {
	AddHeader(header, value string)
	ResetHeaders()
	SetFrom(from string)
	SetSubject(subject string)
}