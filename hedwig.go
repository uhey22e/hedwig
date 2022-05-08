package hedwig

import (
	"bytes"
	"context"
	"errors"
	"net/mail"
)

type Mailer struct {
	d    Driver
	from mail.Address
}

type Config struct {
	From mail.Address
}

type MailerOptionFn func(m *Mailer)

type Driver interface {
	SendMail(ctx context.Context, from mail.Address, to []mail.Address, mail *Mail) error
}

var (
	ErrUnknownDriver = errors.New("unknown driver")
)

// NewMailer is intended for use by drivers only. Do not use in application code.
func NewMailer(d Driver, opts ...MailerOptionFn) *Mailer {
	m := &Mailer{d: d}
	for i := range opts {
		opts[i](m)
	}
	return m
}

// DefaultFrom sets an default email address for From property.
func DefaultFrom(from mail.Address) MailerOptionFn {
	return func(m *Mailer) {
		m.from = from
	}
}

// WithDefaultFrom creates an client with the DefaultFrom option.
func (m *Mailer) WithDefaultFrom(from mail.Address) *Mailer {
	return &Mailer{
		d:    m.d,
		from: from,
	}
}

// SendMail sends an email from the default from address.
func (m *Mailer) SendMail(ctx context.Context, to []mail.Address, mail *Mail) error {
	if m.from.Address == "" {
		return errors.New("default from address is empty")
	}
	return m.d.SendMail(ctx, m.from, to, mail)
}

// SendMail sends an email from specific address.
func (m *Mailer) SendMailFrom(ctx context.Context, from mail.Address, to []mail.Address, mail *Mail) error {
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
