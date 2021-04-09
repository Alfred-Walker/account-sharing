package sqs

import (
	"errors"

	"github.com/Alfred-Walker/AccountSharing/constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	
	
)

// getQueueUrl receives current session and queue name.
// It returns queue url when no error is detected.
func GetQueueUrl(svc *sqs.SQS, queueName *string) (*string, error) {
	if svc == nil {
		return nil, errors.New("SQS: service client cannot be nil.")
	}

	if *queueName == "" {
		return nil, errors.New("SQS: You must supply the name of a queue")
	}

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queueName,
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("SQS: target queue name not found.")
	}

	return result.QueueUrl, nil
}

// Returns messages in the sqs (max 10 messages at once)
func GetMessageList(svc *sqs.SQS, timeout *int64) ([]*sqs.Message, error) {
	queueName := constants.SQS_QUEUE
	
	queueUrl, err := GetQueueUrl(svc, &queueName)
	
	if err != nil {
		return nil, errors.New("SQS: Failed to get queue url")
	}

	// single receive (max 10 messages at once)
	receive, err := ReceiveMsg(svc, queueUrl, timeout)

	if err != nil {
		return nil, errors.New("SQS: Failed to receive messages from the queue.")
	}

	if receive == nil {
		return nil, errors.New("SQS: nil message received.")
	}

	return receive.Messages, nil
}

func ReceiveMsg(svc *sqs.SQS, queueUrl *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	if svc == nil {
		return nil, errors.New("SQS service client is nil.")
	}

	return svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: aws.Int64(constants.SQS_MAX_NUMBER_OF_MESSAGES),
		VisibilityTimeout:   timeout,
	})
}

func DeleteMsg(svc *sqs.SQS, messageHandle *string) error {
	queueName := constants.SQS_QUEUE
	if svc == nil {
		return errors.New("SQS service client is nil.")
	}

	queueUrl, err := GetQueueUrl(svc, &queueName )

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: messageHandle,
	})

	return err
}

func SendMsg(svc *sqs.SQS, queueUrl *string, body *string) error {
	if svc == nil {
		return errors.New("SQS service client is nil.")
	}

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Misc": {
				DataType:    aws.String("String"),
				StringValue: aws.String("None"),
			},
		},
		MessageBody: body,
		QueueUrl:    queueUrl,
	})

	return err
}