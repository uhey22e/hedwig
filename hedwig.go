package hedwig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/mail"
)

type Mailer struct {
	d Driver
}

type Driver interface {
	SendMail(ctx context.Context, from mail.Address, to []mail.Address, mail *Mail) error
}

var (
	drivers          map[string]Driver
	ErrUnknownDriver = errors.New("unknown driver")
)

// NewMailer is intended for use by drivers only. Do not use in application code.
func NewMailer(d Driver) *Mailer {
	return &Mailer{d: d}
}

func OpenMailer(driverName string) (*Mailer, error) {
	drv, ok := drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownDriver, driverName)
	}
	return &Mailer{
		d: drv,
	}, nil
}

// SendMail sends an email.
func (m *Mailer) SendMail(ctx context.Context, from mail.Address, to []mail.Address, mail *Mail) error {
	return m.d.SendMail(ctx, from, to, mail)
}

// Mail is a concrete representation of an email.
// Mail has embedded bytes.Buffer to read/write the message body.
type Mail struct {
	// Subject of the email.
	Subject string
	// Content type of the message body. Defaults to text/plain.
	ContentType ContentType

	// Buffer for the message body.
	bytes.Buffer
}

// Content-Type of the email body.
type ContentType int

const (
	ContentTypePlainText ContentType = iota
	ContentTypeHTML
)

func (c ContentType) String() string {
	return map[ContentType]string{
		ContentTypePlainText: "text/plain",
		ContentTypeHTML:      "text/html",
	}[c]
}
