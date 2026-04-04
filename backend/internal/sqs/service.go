package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ISQSService interface {
	SendMessage(ctx context.Context, message ReplyEmailMessage) error
	SendBatchMessages(ctx context.Context, messages []ReplyEmailMessage) error
}

type SQSService struct {
	client      *sqs.Client
	sqsQueueURL string
}

func NewSQSService(ctx context.Context) (*SQSService, error) {
	cfg := aws.Config{}
	cfg.Region = os.Getenv("AWS_REGION")
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	// Provide fake credentials in test environment
	if os.Getenv("NODE_ENV") == "test" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider("test-key", "test-secret", "")
	} else {
		loadedCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}
		cfg = loadedCfg
	}

	// Comment out the client creation above, and uncomment this to test locally:
	// AWS access and secret key from .env will automatically be used instead of test creds then
	// client := sqs.NewFromConfig(aws.Config{}) // <- Uncomment this to test locally

	return &SQSService{
		client:      sqs.NewFromConfig(cfg),
		sqsQueueURL: os.Getenv("SQS_QUEUE_URL_PROD"),
	}, nil
}

func (s *SQSService) SendMessage(ctx context.Context, message ReplyEmailMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &s.sqsQueueURL,
		MessageBody: aws.String(string(body)),
	})
	return err
}
