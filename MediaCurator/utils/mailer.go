package utils

import (
	"fmt"
	"strings"

	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESClient interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

type Mailer struct {
	from      string
	sesClient SESClient
}

func NewEmailSender(from string) (*Mailer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-2"), // TODO: Make this dynamic
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	sesClient := ses.NewFromConfig(cfg)

	return &Mailer{
		from:      from,
		sesClient: sesClient,
	}, nil
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	formatted := m.formatBody(body)
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: &formatted,
				},
				Text: &types.Content{
					Data: &body,
				},
			},
			Subject: &types.Content{
				Data: &subject,
			},
		},
		Source: &m.from,
	}

	_, err := m.sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (m *Mailer) formatBody(body string) string {
	formatted := strings.ReplaceAll(body, "\n", "<br>")

	return formatted
}