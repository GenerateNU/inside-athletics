package lambdasrc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	sqs "inside-athletics/internal/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var sesService *SESEmailService

func init() {
	var err error
	sesService, err = NewSESEmailService(
		context.Background(),
		getEnv("SES_REGION", "us-east-1"),
		getEnv("SES_FROM_EMAIL", "priseregenerate@gmail.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create SES service: %v", err)
	}
}

func handler(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	log.Printf("Processing %d messages", len(event.Records))

	var batchItemFailures []events.SQSBatchItemFailure

	for _, record := range event.Records {
		err := processRecord(ctx, record)
		if err != nil {
			log.Printf("Failed to process message %s: %v", record.MessageId, err)

			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	log.Printf("Processed %d messages. Failures: %d", len(event.Records), len(batchItemFailures))

	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func processRecord(ctx context.Context, record events.SQSMessage) error {
	var message sqs.ReplyEmailMessage
	if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
		return fmt.Errorf("failed to parse message body: %w", err)
	}

	if message.To == "" || message.From == "" || message.Message == "" {
		return errors.New("missing required fields in message")
	}

	return sesService.SendReplyEmail(ctx, message)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	lambda.Start(handler)
}
