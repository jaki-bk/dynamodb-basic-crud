package controllers

import (
	"dynamodb-basic-crud/models"
	"dynamodb-basic-crud/services"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	fmt.Println(user)

	if err := c.service.CreateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	name := ctx.FormValue("name")
	email := ctx.FormValue("email")

	fmt.Println(name)

	err := c.service.UpdateUser(id, name, email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update user: %v", err)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (c *UserController) GetUserByID(ctx echo.Context) error {
	id := ctx.Param("id")
	user, err := c.service.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	err := c.service.DeleteUser(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete user: %v", err)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (c *UserController) GetAllUsers(ctx echo.Context) error {
	users, err := c.service.GetAllUsers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to fetch users: %v", err)})
	}

	return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) BulkCreateUsers(ctx echo.Context) error {
	var users []models.User
	if err := ctx.Bind(&users); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	err := c.service.BatchCreateUsers(users)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert users"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "bulk insert successfully"})
}

func (c *UserController) GetUsersByCityAndAge(ctx echo.Context) error {
	city := ctx.QueryParam("city")
	age := ctx.QueryParam("age")

	users, err := c.service.GetUsersByCityAndAge(city, age)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users"})
	}

	return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUsersByEmail(ctx echo.Context) error {
	email := ctx.QueryParam("email")

	users, err := c.service.GetUsersByEmail(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users by email"})
	}

	return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUsersByStatusAndCreatedAt(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	createdAt := ctx.QueryParam("created_at")

	users, err := c.service.GetUsersByStatusAndCreatedAt(status, createdAt)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users by status and created_at"})
	}

	return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUsersByName(ctx echo.Context) error {
	name := ctx.QueryParam("name")

	users, err := c.service.GetUsersByName(name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users by name"})
	}

	return ctx.JSON(http.StatusOK, users)
}
