package types

import "net/mail"

type Mail struct {
	From    mail.Address
	To      []mail.Address
	Subject string
	Body    string

	// Content type of body. Defaults to "text/plain".
	ContentType ContentType
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
