package repositories

import (
	"context"
	"dynamodb-basic-crud/config"
	"dynamodb-basic-crud/models"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct{}

func (r *UserRepository) CreateUser(user models.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	fmt.Println(item)

	_, err = config.DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      item,
	})

	return err
}

func (r *UserRepository) UpdateUser(id string, name string, email string) error {
	fmt.Println(name)

	updateExpression := "SET #name = :name, #email = :email"
	expressionAttributeValues := map[string]types.AttributeValue{
		":name":  &types.AttributeValueMemberS{Value: name},
		":email": &types.AttributeValueMemberS{Value: email},
	}

	expressionAttributeNames := map[string]string{
		"#name":  "name",
		"#email": "email",
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err := config.DB.UpdateItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	result, err := config.DB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil || result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	return &user, err
}

func (r *UserRepository) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := config.DB.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	}

	output, err := config.DB.Scan(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan users: %w", err)
	}

	var users []models.User
	for _, item := range output.Items {
		var user models.User
		err := attributevalue.UnmarshalMap(item, &user)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
