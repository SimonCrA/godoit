package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/simoncra/godoit/config"
	"github.com/simoncra/godoit/internal/handlers"
	"github.com/simoncra/godoit/internal/middlewares"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// declaramos el middleware del token jwt
	authMiddleware := middlewares.NewAuthMiddleware(config.AppSecret)

	// main route
	app.Get("/", handlers.HomeHandler)
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/login", handlers.LoginHandler(db))
	v1.Post("/signup", handlers.SignupHandler(db))
	// v1.Get("/users", jwt, handlers.GetUsersHandler(db))

	// creamos un grupo tasks para manejar todas las rutas de tareas
	task := v1.Group("/tasks")

	task.Use(authMiddleware)

	task.Post("/", handlers.AddTaskHandler(db))
	task.Get("/", handlers.ListTasksHandler(db))
	task.Put("/", handlers.UpdateTaskHandler(db))
}
