package handlers

import (
	"errors"
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
		return c.JSON("true")
	}
}

func UpdateTaskHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task models.Task
		var taskDb models.Task

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

		reqQuery := c.Queries()
		idTask := reqQuery["idTask"]

		if err := db.Where("tasks.id = ? and tasks.fk_id_user = ?", idTask, userId).First(&taskDb).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Handle "record not found" error
				return c.Status(fiber.StatusNotFound).SendString("Task not found")
			} else {
				// Handle other errors
				return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving task from the database")
			}
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

func ListTasksHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from the context
		userJwt := c.Locals("user").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		userId := claims["ID"].(float64)

		// idTask := c.Params("idTask")
		// limit, _ := strconv.Atoi(c.Params("limit", "10"))
		// page, _ := strconv.Atoi(c.Params("page", "0"))

		var tasksDb []models.TaskApiResponse

		err := db.Table("tasks").
			// Where("current_task_date <= ?", time.Now()).

			Select("tasks.id, tasks.title, tasks.description, tasks.fk_id_cat_status, tasks.current_task_date, cat_statuses.name as Name").
			Joins("left join users on users.id = tasks.fk_id_user").
			Joins("left join cat_statuses on cat_statuses.id = tasks.fk_id_cat_status").
			Where("users.id = ? AND tasks.logical_delete = ?", userId, false).
			Order("tasks.created_at desc").
			Scan(&tasksDb).
			Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving users from the database")
		}

		todayTasks := []models.TaskApiResponse{}
		tomorrowTasks := []models.TaskApiResponse{}

		currentTime := time.Now()
		tomorrowTime := currentTime.Add(24 * time.Hour)

		for _, task := range tasksDb {
			currentDbTaskDate, err := time.Parse(time.RFC3339Nano, task.CurrentTaskDate)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, parsing dates.")
			}

			if currentDbTaskDate.Day() == currentTime.Day() && currentDbTaskDate.Month() == currentTime.Month() && currentDbTaskDate.Year() == currentTime.Year() {
				todayTasks = append(todayTasks, task)
			} else if currentDbTaskDate.Day() == tomorrowTime.Day() && currentDbTaskDate.Month() == tomorrowTime.Month() && currentDbTaskDate.Year() == tomorrowTime.Year() {
				tomorrowTasks = append(tomorrowTasks, task)
			}
		}

		taskMap := map[string][]models.TaskApiResponse{
			"today":    todayTasks,
			"tomorrow": tomorrowTasks,
		}

		return c.JSON(taskMap)
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
