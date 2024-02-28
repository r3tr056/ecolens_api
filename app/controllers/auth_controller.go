// Package controllers provides API handlers for the application
package controllers

import (
	"time"

	"github.com/r3tr056/ecolens_api/app/models"
	"github.com/r3tr056/ecolens_api/pkg/utils"
	"github.com/r3tr056/ecolens_api/platform/db"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// @Summary User SignUp
// @Description Create a new user account.
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.SignUp true "User SignUp details"
// @Success 200 {object} fiber.Map{"error":false, "message": "User created successfully", "inserted_id": "123", "user": {"id": "123", "created_at": "2022-01-01T12:00:00Z", "email": "user@example.com", "user_status": 1, "user_role": "user"}}
// @Failure 400 {object} fiber.Map{"error":true, "msg": "Bad Request"}
// @Failure 401 {object} fiber.Map{"error":true, "msg": "Unauthorized"}
// @Failure 500 {object} fiber.Map{"error":true, "msg": "Internal Server Error"}
// @Router /signup [post]
func UserSignUp(c *fiber.Ctx) error {
	// create new user auth struct
	signUp := &models.SignUp{}

	if err := c.BodyParser(signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidateErrors(err),
		})
	}

	user := &models.User{}
	// fill up user data
	user.CreatedAt = time.Now()
	user.Email = signUp.Email
	hashedPassword, err := utils.GeneratePassword(signUp.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.PasswordHash = hashedPassword
	user.UserStatus = 1
	user.UserRole = signUp.UserRole

	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidateErrors(err),
		})
	}

	if err := db.PostgresDB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Delete the password hash
	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"error":       false,
		"msg":         nil,
		"inserted_id": user.ID,
		"user":        user,
	})
}

// @Summary User SignIn
// @Description Authenticate a user and generate access tokens.
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.SignIn true "User SignIn details"
// @Success 200 {object} fiber.Map{"error":false, "message": "Login success", "tokens": {"access": "access_token", "refresh": "refresh_token"}}
// @Failure 400 {object} fiber.Map{"error":true, "msg": "Bad Request"}
// @Failure 401 {object} fiber.Map{"error":true, "msg": "Invalid credentials"}
// @Failure 500 {object} fiber.Map{"error":true, "msg": "Internal Server Error"}
// @Router /signin [post]
func UserSignIn(c *fiber.Ctx) error {
	signIn := &models.SignIn{}
	var user models.User

	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := db.PostgresDB.Where("email = ?", signIn.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signIn.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid credentials",
		})
	}

	tokens, err := utils.GenerateNewTokens(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Login success",
		"userId":  user.ID,
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

// func ForgotPassword(c *fiber.Ctx) error {
// 	var request models.ForgotPassword
// 	var user models.User

// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid request payload",
// 		})
// 	}

// 	if err := db.PostgresDB.Where("email = ?", request.Email).First(&user).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   "User not found",
// 		})
// 	}

// }
