package hedwig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type AmazonSESClient struct {
	sesClient *sesv2.Client
}

func NewAmazonSESClient(cfg aws.Config) *AmazonSESClient {
	return &AmazonSESClient{
		sesClient: sesv2.NewFromConfig(cfg),
	}
}

func (c *AmazonSESClient) SendMail(ctx context.Context, m *EMail) error {
	params := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(m.From.String()),
		Destination: &types.Destination{
			ToAddresses: make([]string, len(m.To)),
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: &m.Subject,
				},
			},
		},
	}
	for i, t := range m.To {
		params.Destination.ToAddresses[i] = t.String()
	}
	switch m.ContentType {
	case ContentTypePlainText:
		params.Content.Simple.Body = &types.Body{
			Text: &types.Content{
				Data: &m.Body,
			},
		}
	case ContentTypeHTML:
		params.Content.Simple.Body = &types.Body{
			Html: &types.Content{
				Data: &m.Body,
			},
		}
	}
	_, err := c.sesClient.SendEmail(ctx, params)
	return err
}
