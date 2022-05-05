package hedwig

import (
	"context"
	"net"
	"net/smtp"
	"os"
	"testing"
)

func TestSMTPClient_SendMail(t *testing.T) {
	// Test run without error
	addr, ok := os.LookupEnv("TEST_SMTP_SERVER")
	if !ok {
		t.Skip()
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	c := &SMTPClient{
		Address: addr,
		Auth:    smtp.PlainAuth("", os.Getenv("TEST_SMTP_USERNAME"), os.Getenv("TEST_SMTP_PASSWORD"), host),
	}
	ctx := context.TODO()
	if err := c.SendMail(ctx, content); err != nil {
		t.Errorf("SMTPClient.SendMail() error = %v", err)
	}
}
