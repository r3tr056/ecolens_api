package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/r3tr056/ecolens_api/app/models"
	"github.com/r3tr056/ecolens_api/platform/db"
	"github.com/r3tr056/ecolens_api/platform/pubsub"
)

var SearchTaskRPC *pubsub.PubsubClient

func StartSearchTaskRPC() {
	var err error
	SearchTaskRPC, err = pubsub.NewPubSubClient(
		os.Getenv("GOOGLE_PROJECT_ID"),
		os.Getenv("SEARCH_TOPIC_NAME"),
		os.Getenv("SEARCH_SUB_NAME"),
	)
	if err != nil {
		log.Fatalf("Error starting Pub/Sub Client for Search: %v", err)
	}
}

// @Summary Perform full-text search on product names
// @Description Returns matching product names based on the provided search term
// @ID matchTS
// @Accept json
// @Produce json
// @Param term query string true "Search term for full-text search"
// @Success 200 {array} models.MatchResult
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Router /matchTS [get]
func MatchTS(c *fiber.Ctx) error {
	term := c.Query("term")

	if len(term) > 1 {
		var products []models.Product
		db.PostgresDB.Table("products").Where("ts @@ plainto_tsquery('simple', ?)", fmt.Sprintf("%s:*", term)).Select("name").Limit(10).Find(&products)
		if len(products) > 0 {
			result := make([]models.MatchResult, len(products))
			for i, product := range products {
				result[i] = models.MatchResult{
					Rank:    i,
					Content: product.Name,
				}
			}
			return c.JSON(result)
		} else {
			return c.SendStatus(fiber.StatusNotFound)
		}
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

// @Summary Perform a product search
// @Description Perform a search for products based on the specified search term
// @Tags Products
// @Accept json
// @Produce json
// @Param searchTerm query string true "Search term for products"
// @Param page query integer false "Page number (default is 1)"
// @Param pageSize query integer false "Page size (default is 10)"
// @Success 201 {object} models.SearchResultPage "Successful response with search results"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/products/search [get]
func PerformProductSearch(c *fiber.Ctx) error {
	ctx := context.Background()

	searchTerm := c.FormValue("searchTerm")
	page, err := strconv.Atoi(c.FormValue("page", "1"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.FormValue("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	products, err := models.SearchProducts(ctx, db.RedisClient, db.PostgresDB, searchTerm, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to perform product search",
		})
	}

	var productResults []models.SearchResult
	for i, result := range products {
		productResults = append(productResults, models.CreateSearchResult(i, result.Name, result.Description))
	}

	pageID := models.GenerateRandomPageID()
	searchResultPage := models.SearchResultPage{
		PageID:  pageID,
		Index:   page,
		Length:  pageSize,
		Results: productResults,
	}

	return c.Status(fiber.StatusCreated).JSON(searchResultPage)
}

// @Summary Perform a marketplace product search
// @Description Perform a search for marketplace products based on the specified search term
// @Tags MarketplaceProducts
// @Accept json
// @Produce json
// @Param searchTerm query string true "Search term for marketplace products"
// @Param page query integer false "Page number (default is 1)"
// @Param pageSize query integer false "Page size (default is 10)"
// @Success 201 {object} models.SearchResultPage "Successful response with search results"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/marketplace/products/search [get]
func PerformMarketplaceProductSearch(c *fiber.Ctx) error {
	ctx := context.Background()

	searchTerm := c.FormValue("searchTerm")
	page, err := strconv.Atoi(c.FormValue("page", "1"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.FormValue("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	marketProducts, err := models.SearchMarketplaceProducts(ctx, db.RedisClient, db.PostgresDB, searchTerm, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to perform product search",
		})
	}

	var productResults []models.SearchResult
	for i, result := range marketProducts {
		productResults = append(productResults, models.CreateSearchResult(i, result.Name, result.Description))
	}

	pageID := models.GenerateRandomPageID()
	searchResultPage := models.SearchResultPage{
		PageID:  pageID,
		Index:   page,
		Length:  pageSize,
		Results: productResults,
	}

	return c.Status(fiber.StatusCreated).JSON(searchResultPage)
}

// @Summary Perform a report search
// @Description Perform a search for reports based on the specified search term
// @Tags Reports
// @Accept json
// @Produce json
// @Param searchTerm query string true "Search term for reports"
// @Param page query integer false "Page number (default is 1)"
// @Param pageSize query integer false "Page size (default is 10)"
// @Success 201 {object} models.SearchResultPage "Successful response with search results"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/reports/search [get]
func PerformReportSearch(c *fiber.Ctx) error {
	ctx := context.Background()

	searchTerm := c.FormValue("searchTerm")
	page, err := strconv.Atoi(c.FormValue("page", "1"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.FormValue("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	reports, err := models.SearchReports(ctx, db.RedisClient, db.PostgresDB, searchTerm, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to perform product search",
		})
	}

	var productResults []models.SearchResult
	for i, result := range reports {
		productResults = append(productResults, models.CreateSearchResult(i, result.Name, result.Summary))
	}

	pageID := models.GenerateRandomPageID()
	searchResultPage := models.SearchResultPage{
		PageID:  pageID,
		Index:   page,
		Length:  pageSize,
		Results: productResults,
	}

	return c.Status(fiber.StatusCreated).JSON(searchResultPage)
}

// @Summary Perform an image search
// @Description Perform a search using the submitted image
// @Tags Image Search
// @Accept json
// @Produce json
// @Param userMeta body models.UserMeta true "User metadata including UserID"
// @Param image formData file true "Image file to be searched"
// @Success 201 {object} fiber.Map "Successful response with image search result"
// @Failure 400 {object} ErrorResponse "Bad request, failed to parse JSON data or open image file"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/images/search [post]
func PerformImageSearch(c *fiber.Ctx) error {
	// get the user id
	userMeta := new(models.UserMeta)
	if err := c.BodyParser(userMeta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to parse JSON data",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to open the image file: %v", err),
		})
	}

	imageFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to open the image file: %v", err),
		})
	}
	defer imageFile.Close()

	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read the image file",
		})
	}

	imageURL, err := models.UploadImageToCloudStorage(db.PostgresDB, uint(userMeta.UserID), "application/jpeg", "", false, imageBytes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to upload image",
		})
	}

	args := map[string]interface{}{
		"image_reference": imageURL,
	}
	messageID, err := SearchTaskRPC.PublishMessage("image-search", args)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to call RPC method : %v", err),
		})
	}

	result, err := SearchTaskRPC.WaitForResponse(messageID, 1*time.Second, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to fetch result from RPC method : %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": result})
}
