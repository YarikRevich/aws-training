package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TableName = "university-ladger"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(options *config.LoadOptions) error {
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		_, err = svc.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			TableName: aws.String(TableName),
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeN,
				},
				{
					AttributeName: aws.String("year"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("year"),
					KeyType:       types.KeyTypeRange,
				},
			},
		})

		if err != nil {
			panic(err)
		}
	}

	_, err = svc.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberN{Value: "1"},
			"year": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	if err != nil {
		panic(err)
	}

	out, err := svc.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(TableName),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
