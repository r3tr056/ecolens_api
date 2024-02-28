package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/r3tr056/ecolens_api/app/models"
	"github.com/r3tr056/ecolens_api/platform/db"
)

// pubSubClient, err := pubsub_client.NewPubSubClient("product", "product")

// AddProduct godoc
// @Summary Add a new product
// @Description Adds a new product to the database and triggers analysis of product information.
// @Accept json
// @Produce json
// @Param newProduct body models.Product true "New product information to add"
// @Success 201 {object} models.Product "Product added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request or product data"
// @Failure 500 {object} ErrorResponse "Failed to create product or analyze product information"
// @Router /products [post]
func AddProduct(c *fiber.Ctx) error {
	var newProduct models.Product
	if err := c.BodyParser(&newProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add the new product to the database
	if err := db.PostgresDB.Create(&newProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// _, err = pubSubClient.CallMethod("analyze-product-info", newProduct.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to analyze product information.",
	// 	})
	// }

	return c.Status(fiber.StatusCreated).JSON(newProduct)
}

// AddMarketPlaceProduct godoc
// @Summary Add a new marketplace product
// @Description Adds a new product to the marketplace database and triggers analysis of product information.
// @Accept json
// @Produce json
// @Param newProduct body models.MarketPlaceProduct true "New marketplace product information to add"
// @Success 201 {object} models.MarketPlaceProduct "Marketplace product added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request or product data"
// @Failure 500 {object} ErrorResponse "Failed to create marketplace product or analyze product information"
// @Router /marketplace/products [post]
func AddMarketPlaceProduct(c *fiber.Ctx) error {
	var newProduct models.MarketPlaceProduct
	if err := c.BodyParser(&newProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.PostgresDB.Create(&newProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// pubsubClient, err := pubsub_client.NewPubSubClient("product", "marketplace")
	// _, err = pubsubClient.CallMethod("analyze-product-info", newProduct.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to analyze product information.",
	// 	})
	// }

	return c.Status(fiber.StatusCreated).JSON(newProduct)
}

// UpdateProduct godoc
// @Summary Update an existing product by ID
// @Description Updates an existing product in the database by the specified ID and triggers analysis of product information.
// @Accept json
// @Produce json
// @Param id path integer true "Product ID to update"
// @Param updatedProduct body models.Product true "Updated product information"
// @Success 200 {object} models.Product "Product updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request, product ID, or product data"
// @Failure 500 {object} ErrorResponse "Failed to update product or analyze product information"
// @Router /products/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Product ID",
		})
	}

	var updatedProduct models.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update the existing product in the database
	if err := db.PostgresDB.Model(&models.Product{}).Where("id = ?", id).Updates(updatedProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product.",
		})
	}

	// pubSubClient, err := pubsub_client.NewPubSubClient("product", "product")
	// _, err = pubSubClient.CallMethod("analyze-product-info", updatedProduct.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to analyze product information.",
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(updatedProduct)
}

// GetProducts godoc
// @Summary Get a list of products with pagination
// @Description Retrieves a paginated list of products based on the specified page and limit parameters.
// @Accept json
// @Produce json
// @Param page query integer false "Page number for pagination (default is 1)"
// @Param limit query integer false "Number of products to retrieve per page (default is 10)"
// @Success 200 {array} models.Product "Successful response with the list of products"
// @Failure 400 {object} ErrorResponse "Invalid page or limit parameter"
// @Failure 500 {object} ErrorResponse "Failed to retrieve products"
// @Router /products [get]
func GetProducts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	// Retrieve paginated products from the database
	var products []models.Product
	db.PostgresDB.Offset(offset).Limit(limit).Find(&products)

	return c.JSON(products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieves product details based on the specified ID, including related data such as images, marketplace alternatives, LCAMetrics, and reports.
// @Accept json
// @Produce json
// @Param id path integer true "Product ID to retrieve"
// @Success 200 {object} models.Product "Successful response with the product details"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /products/{id} [get]
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var product models.Product
	if err := db.PostgresDB.Preload("Images").Preload("MarketplaceAlternatives").Preload("EPD.LCAMetrics").Preload("Reports").First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}
