package dtos

type User struct {
	ID   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
	Age  string `json:"age" dynamodbav:"age"`
}
