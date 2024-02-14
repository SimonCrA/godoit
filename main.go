package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/simoncra/godoit/config"
	"github.com/simoncra/godoit/internal/models"
	"github.com/simoncra/godoit/internal/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// initialize the database
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// auto-migrate the database models
	db.AutoMigrate(&models.CatStatus{}, &models.CatCategory{}, &models.User{}, &models.Task{})

	// create the fiber app
	app := fiber.New()

	// set up middlewares
	app.Use(fiberLogger.New())
	app.Use(cors.New())
	app.Use(recover.New())
	// app.Use(middlewares.NotFoundHandler)

	// Set up routes
	routes.SetupRoutes(app, db)

	err = app.Listen(":" + config.AppPort)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
