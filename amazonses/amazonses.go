package amazonses

import (
	"context"

	"github.com/uhey22e/hedwig"
	"github.com/uhey22e/hedwig/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// Client for Amazon SES.
// Using SESv2 API with aws-sdk-go-v2.
type Mailer struct {
	client *sesv2.Client
}

func NewAmazonSESClient(cfg aws.Config) *hedwig.Mailer {
	d := &Mailer{
		client: sesv2.NewFromConfig(cfg),
	}
	return hedwig.NewMailer(d)
}

func (c *Mailer) SendMail(ctx context.Context, m *types.Mail) error {
	params := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(m.From.String()),
		Destination: &sesTypes.Destination{
			ToAddresses: make([]string, len(m.To)),
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Subject: &sesTypes.Content{
					Data: &m.Subject,
				},
			},
		},
	}
	for i, t := range m.To {
		params.Destination.ToAddresses[i] = t.String()
	}
	switch m.ContentType {
	case types.ContentTypePlainText:
		params.Content.Simple.Body = &sesTypes.Body{
			Text: &sesTypes.Content{
				Data: &m.Body,
			},
		}
	case types.ContentTypeHTML:
		params.Content.Simple.Body = &sesTypes.Body{
			Html: &sesTypes.Content{
				Data: &m.Body,
			},
		}
	}
	_, err := c.client.SendEmail(ctx, params)
	return err
}
