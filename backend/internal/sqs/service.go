package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ISQSService interface {
	SendMessage(ctx context.Context, message DisasterEmailMessage) error
	SendBatchMessages(ctx context.Context, messages []DisasterEmailMessage) error
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

func (s *SQSService) SendMessage(ctx context.Context, message DisasterEmailMessage) error {
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

func (s *SQSService) SendBatchMessages(ctx context.Context, messages []DisasterEmailMessage) error {
	if len(messages) == 0 {
		return nil
	}

	log.Printf("Got %d messages to SQS Batch send", len(messages))

	for i := 0; i < len(messages); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(messages) {
			end = len(messages)
		}
		batch := messages[i:end]

		entries := make([]sqstypes.SendMessageBatchRequestEntry, len(batch))
		for index, message := range batch {
			body, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}
			id := fmt.Sprintf("%d", i+index)
			entries[index] = sqstypes.SendMessageBatchRequestEntry{
				Id:          &id,
				MessageBody: aws.String(string(body)),
			}
		}

		log.Printf("Entries: \n%v", entries)

		command := &sqs.SendMessageBatchInput{
			QueueUrl: &s.sqsQueueURL,
			Entries:  entries,
		}

		prettyMessages, _ := json.MarshalIndent(messages, "", "  ")
		log.Printf("COMMAND: \n\n\n%s\n\n\n", prettyMessages)

		response, err := s.client.SendMessageBatch(ctx, command)
		if err != nil {
			log.Printf("Error sending batch messages: %v", err)
			return err
		}

		prettyResponse, _ := json.MarshalIndent(response, "", "  ")
		log.Printf("Sending batch messages response: %s", prettyResponse)
		log.Printf("Response failed: %v", response.Failed)
		log.Printf("Response successful: %v", response.Successful)

		// Log successful and failed messages
		if len(response.Successful) > 0 {
			log.Printf("Successfully sent %d messages", len(response.Successful))
		}
		if len(response.Failed) > 0 {
			log.Printf("Failed to send %d messages: %v", len(response.Failed), response.Failed)
		}
	}

	return nil
}
