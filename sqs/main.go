package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
)

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const QueueName = "yaroslav"

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), func(options *config.LoadOptions) error {
		return nil
	})

	if err != nil {
		panic(err)
	}

	svc := sqs.NewFromConfig(cfg)

	result, err := svc.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: aws.String(QueueName),
	})

	if err != nil {
		panic(err)
	}

	_, err = svc.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    result.QueueUrl,
		MessageBody: aws.String("IT WORKS"),
	})

	if err != nil {
		panic(err)
	}
}
