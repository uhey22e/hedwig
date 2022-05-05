package hedwig

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"net/smtp"
)

type SMTPClient struct {
	Address string
	Auth    smtp.Auth
}

const (
	initialBufferSize = 1024
)

var (
	mimeEncoder = mime.BEncoding
	crlf        = []byte("\r\n")
)

func (c *SMTPClient) SendMail(ctx context.Context, m *Mail) error {
	buf := &bytes.Buffer{}
	buf.Grow(initialBufferSize + len(m.Body)*4/3)
	err := writeSMTPHeader(buf, m)
	if err != nil {
		return err
	}
	_, err = buf.Write(crlf)
	if err != nil {
		return err
	}

	enc := NewSMTPBodyEncoder(buf)
	defer enc.Close()
	_, err = enc.WriteString(m.Body)
	if err != nil {
		return err
	}
	err = enc.Close()
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

func writeSMTPHeader(w io.Writer, m *Mail) (err error) {
	lines := []string{
		"Content-Transfer-Encoding: base64",
		"Content-Type: " + m.ContentType.Value(),
		"From: " + mimeEncoder.Encode(CharSet, formatAddress(m.From)),
		"To: " + mimeEncoder.Encode(CharSet, formatAddresses(m.To)),
		"Bcc: " + mimeEncoder.Encode(CharSet, formatAddresses(m.BCC)),
		"Subject: " + mimeEncoder.Encode(CharSet, m.Subject),
	}
	for _, l := range lines {
		_, err = io.WriteString(w, l+"\r\n")
		if err != nil {
			return
		}
	}
	return nil
}

type SMTPBodyEncoder struct {
	dest io.Writer
	b64  io.WriteCloser
}

func NewSMTPBodyEncoder(w io.Writer) *SMTPBodyEncoder {
	return &SMTPBodyEncoder{
		dest: w,
		b64:  base64.NewEncoder(base64.RawStdEncoding, w),
	}
}

func (e *SMTPBodyEncoder) WriteString(s string) (int, error) {
	l := ((78 - 2) * 3) / 4
	b := []byte(s)
	sum := 0
	for {
		chunk := b
		if r := len(chunk); r == 0 {
			break
		} else if r > l {
			chunk = chunk[:l]
		}
		n, err := e.b64.Write(chunk)
		if err != nil {
			return n, err
		}
		m, err := e.dest.Write(crlf)
		if err != nil {
			return n, err
		}
		b = b[n:]
		sum += (n + m)
	}
	return sum, nil
}

func (e *SMTPBodyEncoder) Close() error {
	return e.b64.Close()
}
