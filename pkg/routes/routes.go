package routes

import (
	"github.com/r3tr056/ecolens_api/app/controllers"
	"github.com/r3tr056/ecolens_api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-TOKEN",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.Use(middleware.HeartBeat("/ping"))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/user/signin", controllers.UserSignIn)
	v1.Post("/user/signup", controllers.UserSignUp)

	// List all private routes
	// user routes
	v1.Get("/user/users", middleware.JWTProtected(), controllers.GetUsersHandler)
	v1.Get("/user", middleware.JWTProtected(), controllers.GetUserHandler)
	v1.Put("/user", middleware.JWTProtected(), controllers.UpdateUserHandler)
	v1.Delete("/user", middleware.JWTProtected(), controllers.DeleteUserHandler)

	// search routes
	v1.Post("/autocomplete", middleware.JWTProtected(), controllers.MatchTS)

	// report search
	v1.Post("/report/search", middleware.JWTProtected(), controllers.PerformReportSearch)

	// product routes
	v1.Post("/product/search", middleware.JWTProtected(), controllers.PerformProductSearch)
	v1.Post("/product", middleware.JWTProtected(), controllers.AddProduct)
	v1.Get("/product", middleware.JWTProtected(), controllers.GetProductByID)
	v1.Post("/mkplcproduct", middleware.JWTProtected(), controllers.AddMarketPlaceProduct)
	v1.Post("/mkplcproduct/search", middleware.JWTProtected(), controllers.PerformMarketplaceProductSearch)
	v1.Put("/product", middleware.JWTProtected(), controllers.UpdateProduct)
	v1.Get("/products", middleware.JWTProtected(), controllers.GetProducts)

}
