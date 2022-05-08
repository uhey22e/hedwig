package generalsmtp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"net/smtp"
	"strings"
	"unicode/utf8"

	"github.com/uhey22e/hedwig"
	"github.com/uhey22e/hedwig/types"
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

func NewMailer(ctx context.Context, addr string, auth smtp.Auth) (*hedwig.Mailer, error) {
	d := &Mailer{
		address: addr,
		auth:    auth,
	}
	return hedwig.NewMailer(d), nil
}

func (c *Mailer) SendMail(ctx context.Context, m *types.Mail) error {
	buf := &bytes.Buffer{}
	buf.Grow(initialBufLen)
	err := writeSMTPHeaders(buf, m)
	if err != nil {
		return err
	}
	_, err = buf.Write(crlf)
	if err != nil {
		return err
	}
	err = writeSMTPBody(buf, m.Body)
	if err != nil {
		return err
	}
	to := make([]string, len(m.To))
	for i, t := range m.To {
		to[i] = t.Address
	}
	return smtp.SendMail(c.address, c.auth, m.From.Address, to, buf.Bytes())
}

func writeSMTPHeaders(w io.Writer, m *types.Mail) (err error) {
	lines := []string{
		"From: " + m.From.String(),
		"To: " + encodeAddresses(m.To),
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
	return nil
}

func writeSMTPBody(w io.Writer, body string) error {
	enc := base64.NewEncoder(base64.StdEncoding, w)
	defer enc.Close()

	// If the content is short, do not bother splitting the encoded-word.
	if base64.StdEncoding.EncodedLen(len(body)) <= maxEncodedWordLen {
		_, err := io.WriteString(enc, body)
		return err
	}

	var err error
	var currentLen, last, runeLen int
	for i := 0; i < len(body); i += runeLen {
		// Multi-byte characters must not be split across encoded-words.
		// See RFC 2047, section 5.3.
		_, runeLen = utf8.DecodeRuneInString(body[i:])
		if currentLen+runeLen <= maxBase64Len {
			currentLen += runeLen
		} else {
			_, err = io.WriteString(enc, body[last:i])
			if err != nil {
				return err
			}
			_, err = io.WriteString(w, "\r\n")
			if err != nil {
				return err
			}
			last = i
			currentLen = runeLen
		}
	}
	_, err = io.WriteString(enc, body[last:])
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
