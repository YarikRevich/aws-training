package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"os"
	"time"
)

var USER_UUID = uuid.New()

func prepareChat(chatName string, svc *sqs.Client) *string {
	result, err := svc.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(chatName),
	})

	if err != nil {

		result, err := svc.CreateQueue(context.Background(), &sqs.CreateQueueInput{
			QueueName: aws.String(chatName),
			Attributes: map[string]string{
				"DelaySeconds":      "0",
				"VisibilityTimeout": "0",
			},
		})

		if err != nil {
			panic(err)
		}

		return result.QueueUrl
	}
	return result.QueueUrl
}

func receiveMessages(chatUrl *string, svc *sqs.Client) {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		result, err := svc.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:              chatUrl,
			MessageAttributeNames: []string{"id"},
		})

		if err != nil {
			panic(err)
		}

		for _, m := range result.Messages {
			v, ok := m.MessageAttributes["id"]
			if !ok {
				continue
			}

			if *v.StringValue == USER_UUID.String() {
				continue
			}

			fmt.Println(*m.Body)

			_, err := svc.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      chatUrl,
				ReceiptHandle: m.ReceiptHandle,
			})

			if err != nil {
				panic(err)
			}
		}

	}
}

func writeMessage(message string, chatUrl *string, svc *sqs.Client) {
	_, err := svc.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    chatUrl,
		MessageBody: aws.String(message),

		MessageAttributes: map[string]types.MessageAttributeValue{
			"id": types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(USER_UUID.String()),
			},
		},
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), func(options *config.LoadOptions) error {
		return nil
	})

	if err != nil {
		panic(err)
	}

	svc := sqs.NewFromConfig(cfg)

	fmt.Println("Write chat name to connect to")

	scanner := bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var chatUrl *string
	for scanner.Scan() {
		chatName := scanner.Text()
		fmt.Println("Chosen chat name is: ", chatName)
		chatUrl = prepareChat(chatName, svc)
		break
	}

	go receiveMessages(chatUrl, svc)

	scanner = bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for scanner.Scan() {
		message := scanner.Text()
		writeMessage(message, chatUrl, svc)
	}
}
