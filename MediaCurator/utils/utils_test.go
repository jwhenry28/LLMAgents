package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// MockSESClient is a mock implementation of the SES client for testing
type MockSESClient struct {
	SendEmailFunc func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

func (m *MockSESClient) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	return m.SendEmailFunc(ctx, params, optFns...)
}

func TestMailer_SendEmail(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		subject string
		body    string
		mockFn  func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
		wantErr bool
	}{
		{
			name:    "successful email send",
			from:    "test@example.com",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "Test Body <br><br>",
			mockFn: func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
				// Verify input parameters
				if *params.Source != "test@example.com" {
					t.Errorf("unexpected from address: got %v, want %v", *params.Source, "test@example.com")
				}
				if params.Destination.ToAddresses[0] != "recipient@example.com" {
					t.Errorf("unexpected to address: got %v, want %v", params.Destination.ToAddresses[0], "recipient@example.com")
				}
				if *params.Message.Subject.Data != "Test Subject" {
					t.Errorf("unexpected subject: got %v, want %v", *params.Message.Subject.Data, "Test Subject")
				}
				if *params.Message.Body.Html.Data != "Test Body <br><br>" {
					t.Errorf("unexpected body: got %v, want %v", *params.Message.Body.Html.Data, "Test Body <br><br>")
				}
				return &ses.SendEmailOutput{}, nil
			},
			wantErr: false,
		},
		{
			name:    "failed email send",
			from:    "test@example.com",
			to:      "recipient@example.com",
			subject: "Test Subject",
			body:    "Test Body",
			mockFn: func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
				return nil, fmt.Errorf("mock SES error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockSESClient{
				SendEmailFunc: tt.mockFn,
			}

			m := &Mailer{
				from:      tt.from,
				sesClient: mockClient,
			}

			err := m.SendEmail(tt.to, tt.subject, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mailer.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
