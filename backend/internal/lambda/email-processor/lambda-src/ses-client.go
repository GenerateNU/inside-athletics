package lambdasrc

import (
	"context"
	"os"

	lambda "inside-athletics/internal/lambda/email-processor/emails"
	sqs "inside-athletics/internal/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESEmailService struct {
	client    *ses.Client
	fromEmail string
}

func NewSESEmailService(ctx context.Context, region string, fromEmail string) (*SESEmailService, error) {
	if region == "" {
		region = "us-east-1"
	}

	cfg := aws.Config{}
	cfg.Region = os.Getenv("AWS_REGION")
	if cfg.Region == "" {
		cfg.Region = region
	}

	// Provide fake credentials in test environment
	if os.Getenv("NODE_ENV") == "test" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider("test-key", "test-secret", "")
	} else {
		loadedCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
		if err != nil {
			return nil, err
		}
		cfg = loadedCfg
	}

	// client := ses.NewFromConfig(aws.Config{}) // <- USE THIS FOR REAL RUN - AWS will grab config from OS env vars
	client := ses.NewFromConfig(cfg)

	return &SESEmailService{
		client:    client,
		fromEmail: fromEmail,
	}, nil
}

func (s *SESEmailService) SendReplyEmail(ctx context.Context, message sqs.ReplyEmailMessage) error {
	// Render the email component to HTML and plain text
	htmlBody := lambda.RenderReplyEmailHTML(message)
	textBody := lambda.RenderReplyEmailText(message)

	destination := &types.Destination{
		ToAddresses: []string{message.To},
	}

	command := &ses.SendEmailInput{
		Source:      &s.fromEmail,
		Destination: destination,
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String("New Reply in Inside Athletics"),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(htmlBody),
					Charset: aws.String("UTF-8"),
				},
				Text: &types.Content{
					Data:    aws.String(textBody),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	_, err := s.client.SendEmail(ctx, command)
	return err
}
