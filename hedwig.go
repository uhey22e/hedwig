package hedwig

import (
	"context"
	"net/mail"
	"strings"
)

var (
	CharSet = "utf-8"
)

// Content-Type of the email body.
type ContentType int

const (
	ContentTypePlainText ContentType = iota
	ContentTypeHTML
)

// Get a value for Content-Type header.
// e.g. text/plain; charset="utf-8"
func (t ContentType) Value() string {
	return map[ContentType]string{
		ContentTypePlainText: `text/plain; charset="` + CharSet + `"`,
		ContentTypeHTML:      `text/html; charset="` + CharSet + `"`,
	}[t]
}

type Mail struct {
	From    mail.Address
	To      []mail.Address
	Subject string
	Body    string

	// Content type of body. Defaults to "text/plain".
	ContentType ContentType
	// CC addresses. Defaults to empty.
	CC []mail.Address
	// BCC addresses. Defaults to empty.
	BCC []mail.Address
}

type Client interface {
	// Send an email.
	SendMail(context.Context, *Mail) error
}

func formatAddresses(addrs []mail.Address) string {
	s := &strings.Builder{}
	for i, addr := range addrs {
		if i != 0 {
			s.WriteString(",")
		}
		s.WriteString(formatAddress(addr))
	}
	return s.String()
}

func formatAddress(addr mail.Address) string {
	name := ""
	if addr.Name != "" {
		name = `"` + addr.Name + `" `
	}
	return name + "<" + addr.Address + ">"
}
