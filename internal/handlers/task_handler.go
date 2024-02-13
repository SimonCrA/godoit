package handlers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/simoncra/godoit/internal/models"
	"gorm.io/gorm"
)

func AddTaskHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newTask models.Task

		// Get the user from the context
		userJwt := c.Locals("user").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		userId := claims["ID"].(float64)

		if err := c.BodyParser(&newTask); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		fmt.Println("Welcome 👋" + email)
		fmt.Println(int(userId))

		if err := validateTaskInput(&newTask); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		newTask.FkIdUser = int(userId)
		newTask.FkIdCatStatus = 1
		newTask.FkIdCatCategory = 1

		result := db.Create(&newTask)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating a new user")
		}

		// Return the token
		return c.JSON(models.TaskResponse{
			IdTask: newTask.ID,
		})
	}
}

func UpdateTaskHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task models.Task

		// Get the user from the context
		userJwt := c.Locals("user").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		userId := claims["ID"].(float64)

		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := validateTaskInput(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		idTask := c.Params("idTask")

		fmt.Println(idTask)

		var taskDb models.Task
		if err := db.First(&taskDb, idTask).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving users from the database")
		}

		if !task.CurrentTaskDate.IsZero() {
			taskDb.CurrentTaskDate = task.CurrentTaskDate
		}

		if task.Title != "" {
			taskDb.Title = task.Title
		}

		if task.Description != "" {
			taskDb.Description = task.Description
		}

		if task.CatCategory.ID < 1 {
			taskDb.Description = task.Description
		}

		taskDb.FkIdUser = int(userId)
		taskDb.LastStatusChange = time.Now()

		db.Save(taskDb)
		return c.JSON(taskDb)
	}
}

func validateTaskInput(task *models.Task) error {
	validate := validator.New()
	// validate.RegisterValidation("ignorefields", func(fl validator.FieldLevel) bool {
	// 	return true
	// })
	//
	// validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	// 	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// 	if name == "-" {
	// 		return ""
	// 	}
	// 	return name
	// })

	// Validate the Task struct fields
	err := validate.StructExcept(task, "User", "CatStatus", "CatCategory")
	if err != nil {
		var validationErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation error: %s", validationErrors)
	}

	return nil
}
