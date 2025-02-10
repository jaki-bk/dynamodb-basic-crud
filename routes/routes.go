package routes

import (
	"dynamodb-basic-crud/controllers"
	"dynamodb-basic-crud/repositories"
	"dynamodb-basic-crud/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	userRepo := &repositories.UserRepository{}
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	e.POST("/users", userController.CreateUser)
	e.GET("/users/:id", userController.GetUserByID)
	e.DELETE("/users/:id", userController.DeleteUser)
	e.GET("/users", userController.GetAllUsers)
	e.PUT("/users/:id", userController.UpdateUser)
	e.POST("/bulk-create", userController.BulkCreateUsers)
	e.GET("/users/city-age", userController.GetUsersByCityAndAge)
	e.GET("/users/email", userController.GetUsersByEmail)
	e.GET("/users/status-created-at", userController.GetUsersByStatusAndCreatedAt)
	e.GET("/users/name", userController.GetUsersByName)
}
