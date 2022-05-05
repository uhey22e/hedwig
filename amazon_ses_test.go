package hedwig

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
)

func TestAmazonSESClient_SendMail(t *testing.T) {
	// Test run without error
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	c := NewAmazonSESClient(cfg)
	if err := c.SendMail(ctx, email); err != nil {
		t.Errorf("AmazonSESClient.SendMail() error = %v", err)
	}
}
