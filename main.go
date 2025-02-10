package main

import (
	"dynamodb-basic-crud/config"
	"dynamodb-basic-crud/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitDynamoDB()
	//config.DeleteUserTable()
	config.CreateUserTable()

	e := echo.New()
	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
