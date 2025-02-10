package models

type User struct {
	ID        string `json:"id" dynamodbav:"id"`
	Name      string `json:"name" dynamodbav:"name"`
	Email     string `json:"email" dynamodbav:"email"`
	City      string `json:"city" dynamodbav:"city"`
	Status    string `json:"status" dynamodbav:"status"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
	Age       int    `json:"age" dynamodbav:"age"`
}
