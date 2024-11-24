package utils

import (
	"fmt"

	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Mailer struct {
	from      string
	sesClient *ses.Client
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
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
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
