package driver

import (
	"context"

	"github.com/uhey22e/hedwig/types"
)

type Mailer interface {
	SendMail(ctx context.Context, mail *types.Mail) error
}
