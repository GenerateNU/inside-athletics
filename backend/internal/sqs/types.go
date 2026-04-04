package sqs

const BATCH_SIZE = 10 // SQS max batch size; replace with your BATCH_SIZE constant

type DisasterEmailMessage struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}
