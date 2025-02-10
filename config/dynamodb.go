package config

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DB *dynamodb.Client

func InitDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "dummy",
					SecretAccessKey: "dummy",
					SessionToken:    "dummy",
					Source:          "HardcodedCredentials",
				}, nil
			})),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
	)

	if err != nil {
		log.Fatalf("Failed to load DynamoDB config: %v", err)
	}

	DB = dynamodb.NewFromConfig(cfg)
}

func CreateUserTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("Users"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	_, err := DB.CreateTable(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("✅ Users table created successfully")
	return nil
}

func DeleteUserTable() error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("Users"),
	}

	_, err := DB.DeleteTable(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	fmt.Println("✅ Users table deleted successfully")
	return nil
}
