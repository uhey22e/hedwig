package generalsmtp

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"net/smtp"
	"strings"
	"unicode/utf8"

	"github.com/uhey22e/hedwig"
)

type Mailer struct {
	address string
	auth    smtp.Auth
}

const (
	defaultCharSet = "utf-8"
	initialBufLen  = 1024
	// The maximum length of an encoded-word is 75 characters.
	// See RFC 2047, section 2.
	maxEncodedWordLen = 75
)

var (
	crlf         = []byte("\r\n")
	maxBase64Len = base64.StdEncoding.DecodedLen(maxEncodedWordLen)
)

func OpenMailer(ctx context.Context, addr string, auth smtp.Auth, opts ...hedwig.MailerOptionFn) (*hedwig.Mailer, error) {
	d := &Mailer{
		address: addr,
		auth:    auth,
	}
	return hedwig.NewMailer(d, opts...), nil
}

func (c *Mailer) SendMail(ctx context.Context, from mail.Address, to []mail.Address, m *hedwig.Mail) error {
	msg := &bytes.Buffer{}
	msg.Grow(initialBufLen + base64.StdEncoding.EncodedLen(m.Len()))
	err := writeSMTPHeaders(msg, from, to, m)
	if err != nil {
		return err
	}
	err = writeSMTPBody(msg, m.Bytes())
	if err != nil {
		return err
	}
	tos := make([]string, len(to))
	for i, t := range to {
		tos[i] = t.Address
	}
	return smtp.SendMail(c.address, c.auth, from.Address, tos, msg.Bytes())
}

func writeSMTPHeaders(w io.Writer, from mail.Address, to []mail.Address, m *hedwig.Mail) (err error) {
	lines := []string{
		"From: " + from.String(),
		"To: " + encodeAddresses(to),
		"Subject: " + mime.BEncoding.Encode(defaultCharSet, m.Subject),
		"Content-Type: " + fmt.Sprintf("%s; charset=\"%s\"", m.ContentType, defaultCharSet),
		"Content-Transfer-Encoding: base64",
	}
	for _, l := range lines {
		_, err = io.WriteString(w, l+"\r\n")
		if err != nil {
			return
		}
	}
	_, err = w.Write(crlf)
	if err != nil {
		return err
	}
	return nil
}

func writeSMTPBody(w io.Writer, body []byte) error {
	if !utf8.Valid(body) {
		return errors.New("body must be valid utf-8 string")
	}
	enc := base64.NewEncoder(base64.StdEncoding, w)
	defer enc.Close()

	// If the content is short, do not bother splitting the encoded-word.
	if base64.StdEncoding.EncodedLen(len(body)) <= maxEncodedWordLen {
		_, err := enc.Write(body)
		return err
	}

	var err error
	var currentLen, last, runeLen int
	for i := 0; i < len(body); i += runeLen {
		// Multi-byte characters must not be split across encoded-words.
		// See RFC 2047, section 5.3.
		_, runeLen = utf8.DecodeRune(body[i:])
		if currentLen+runeLen <= maxBase64Len {
			currentLen += runeLen
		} else {
			_, err = enc.Write(body[last:i])
			if err != nil {
				return err
			}
			_, err = w.Write(crlf)
			if err != nil {
				return err
			}
			last = i
			currentLen = runeLen
		}
	}
	_, err = enc.Write(body[last:])
	return err
}

func encodeAddresses(addrs []mail.Address) string {
	s := &strings.Builder{}
	for i, addr := range addrs {
		if i != 0 {
			s.WriteString(",")
		}
		s.WriteString(addr.String())
	}
	return s.String()
}
