package hedwig

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"net/smtp"
	"unicode/utf8"
)

type SMTPClient struct {
	Address string
	Auth    smtp.Auth
}

const (
	initialBufferSize = 1024
	// The maximum length of an encoded-word is 75 characters.
	// See RFC 2047, section 2.
	maxEncodedWordLen = 75
)

var (
	crlf         = []byte("\r\n")
	maxBase64Len = base64.StdEncoding.DecodedLen(maxEncodedWordLen)
)

func (c *SMTPClient) SendMail(ctx context.Context, m *EMail) error {
	buf := &bytes.Buffer{}
	buf.Grow(initialBufferSize + base64.StdEncoding.EncodedLen(len(m.Body)))
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
	for i, s := range m.To {
		to[i] = s.Address
	}
	err = smtp.SendMail(c.Address, c.Auth, m.From.Address, to, buf.Bytes())
	return err
}

func writeSMTPHeaders(w io.Writer, m *EMail) (err error) {
	lines := []string{
		"From: " + m.From.String(),
		"To: " + encodeAddresses(m.To),
		"Subject: " + mime.BEncoding.Encode(CharSet, m.Subject),
		"Content-Type: " + m.ContentType.Value(),
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
