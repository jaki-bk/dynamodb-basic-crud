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

func (r *UserRepository) BatchCreateUsers(users []models.User) error {
	writeRequests := make([]types.WriteRequest, len(users))

	for i, user := range users {
		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"id":         &types.AttributeValueMemberS{Value: user.ID},
					"name":       &types.AttributeValueMemberS{Value: user.Name},
					"email":      &types.AttributeValueMemberS{Value: user.Email},
					"city":       &types.AttributeValueMemberS{Value: user.City},
					"status":     &types.AttributeValueMemberS{Value: user.Status},
					"created_at": &types.AttributeValueMemberS{Value: user.CreatedAt},
					"age":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.Age)},
				},
			},
		}
	}

	_, err := config.DB.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"Users": writeRequests,
		},
	})

	return err
}

func (ur *UserRepository) GetUsersByEmail(email string) ([]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String("GSI_Email"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := config.DB.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetUsersByCityAndAge(city string, age string) ([]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String("GSI_City_Age"),
		KeyConditionExpression: aws.String("city = :city and age = :age"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":city": &types.AttributeValueMemberS{Value: city},
			":age":  &types.AttributeValueMemberN{Value: age},
		},
	}

	result, err := config.DB.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetUsersByStatusAndCreatedAt(status string, createdAt string) ([]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String("GSI_Status_CreatedAt"),
		KeyConditionExpression: aws.String("#status = :status and #created_at = :created_at"),
		ExpressionAttributeNames: map[string]string{
			"#status":     "status",
			"#created_at": "created_at",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: status},
			":created_at": &types.AttributeValueMemberS{Value: createdAt},
		},
	}

	result, err := config.DB.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetUsersByName(name string) ([]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String("GSI_Name"),
		KeyConditionExpression: aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
	}

	result, err := config.DB.Query(context.TODO(), input)
	if err != nil {
		fmt.Println("Error querying DynamoDB:", err)
		return nil, err
	}

	var users []models.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
