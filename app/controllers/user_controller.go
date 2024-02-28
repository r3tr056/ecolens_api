package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/r3tr056/ecolens_api/app/models"
	"github.com/r3tr056/ecolens_api/platform/db"

	"gorm.io/gorm"
)

// GetUsersHandler godoc
// @Summary Get a list of users with pagination
// @Description Retrieves a paginated list of users based on the specified page and limit parameters.
// @Accept json
// @Produce json
// @Param page query integer false "Page number for pagination (default is 1)"
// @Param limit query integer false "Number of users to retrieve per page (default is 10)"
// @Success 200 {array} models.User "Successful response with the list of users"
// @Failure 400 {object} ErrorResponse "Invalid page or limit parameter"
// @Failure 500 {object} ErrorResponse "Failed to retrieve users"
// @Router /users [get]
func GetUsersHandler(c *fiber.Ctx) error {
	defaultPage := 1
	defaultLimit := 10

	page, err := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	offset := (page - 1) * limit
	var users []models.User

	result := db.PostgresDB.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retreive users",
		})
	}

	return c.JSON(users)
}

// GetUserHandler godoc
// @Summary Get a user by ID
// @Description Retrieves a user by the specified ID, including related data such as uploaded images and search history.
// @Accept json
// @Produce json
// @Param id path string true "User ID to retrieve"
// @Success 200 {object} models.User "Successful response with the user details"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to retrieve user"
// @Router /users/{id} [get]
func GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	result := db.PostgresDB.Preload("UploadedImages").Preload("SearchHIstory").First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "User  ot found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retreive user",
		})
	}

	return c.JSON(user)
}

// UpdateUserHandler godoc
// @Summary Update a user's information
// @Description Updates a user's information, including the avatar if provided.
// @Accept json
// @Produce json
// @Param id path string true "User ID to update"
// @Param avatar formData file false "New avatar image for the user"
// @Param updatedUser body models.UserUpdate true "Updated user information"
// @Success 200 "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request or update data"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /users/{id} [put]
func UpdateUserHandler(c *fiber.Ctx) error {
	stringUserID := c.Params("id")
	userID, err := strconv.ParseUint(stringUserID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to parse userID",
		})
	}

	avatarImage, err := c.FormFile("avatar")
	var avatarURL string

	if err == nil {
		file, err := avatarImage.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to open avatar image",
			})
		}
		defer file.Close()

		avatarURL, err = models.UploadAvatar(uint(userID), file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to upload avatar to GCS",
			})
		}
	}

	var updatedUser models.UserUpdate
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	updatedUser.AvatarURL = avatarURL

	var existingUser models.User
	result := db.PostgresDB.Where("id = ?", userID).First(&existingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update user",
		})
	}

	db.PostgresDB.Model(&existingUser).Updates(updatedUser)

	return c.SendStatus(fiber.StatusOK)
}

// DeleteUserHandler godoc
// @Summary Delete a user by ID
// @Description Deletes a user based on the specified ID.
// @Accept json
// @Produce json
// @Param id path string true "User ID to delete"
// @Success 200 "User deleted successfully"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to delete user"
// @Router /users/{id} [delete]
func DeleteUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User

	result := db.PostgresDB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete user",
		})
	}

	// Delete the user
	db.PostgresDB.Delete(user)

	return c.SendStatus(fiber.StatusOK)
}
