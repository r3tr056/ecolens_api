package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/r3tr056/ecolens_api/app/controllers"
	"github.com/r3tr056/ecolens_api/pkg/middleware"
	"github.com/r3tr056/ecolens_api/pkg/routes"
	"github.com/r3tr056/ecolens_api/platform/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// open the connection to postgres database
	db.OpenPostgresConnection()
	db.CustomMigrate()

	app := fiber.New()

	// Register middlewares
	middleware.FiberMiddleware(app)

	// Start RPC clients
	controllers.StartSearchTaskRPC()

	// TODO : Routes
	routes.SetupRoutes(app)
	routes.SwaggerRoute(app)
	routes.NotFoundRoute(app)

	// start server
	if os.Getenv("STAGE_STATUS") == "dev" {
		// Server run without gracefull shutdown
		fiberUrl := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
		if err := app.Listen(fiberUrl); err != nil {
			log.Printf("Oops... Server is not running! Reason : %v", err)
		}
	} else {
		// Server run with graceful shutdown
		// create channels for idle conenctions
		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// Received interrupt. Shutdown
			if err := app.Shutdown(); err != nil {
				log.Printf("Oops.... Server is not shutting down! Reason : %v", err)
			}

			close(idleConnsClosed)
		}()

		// Run Server
		fiberUrl := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
		if err := app.Listen(fiberUrl); err != nil {
			log.Printf("Oops... Server is not running! Reason : %v", err)
		}

		<-idleConnsClosed
	}
}
