package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
)

const TemplateName = "test-template"

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), func(options *config.LoadOptions) error {
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := pinpoint.NewFromConfig(cfg)

	_, err = svc.SendMessages(context.Background(), &pinpoint.SendMessagesInput{

		ApplicationId: aws.String("e2e2cd0637bd49eeac9df812db7a613c"),
		MessageRequest: &types.MessageRequest{
			Addresses: map[string]types.AddressConfiguration{
				"+48451128170": types.AddressConfiguration{
					ChannelType: types.ChannelTypeSms,
				},
			},
			MessageConfiguration: &types.DirectMessageConfiguration{
				SMSMessage: &types.SMSMessage{
					Body:        aws.String("Hello, it works, you know!"),
					MessageType: types.MessageTypePromotional,
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}
}
