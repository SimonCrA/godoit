package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/simoncra/godoit/internal/models"
	"gorm.io/gorm"
)

func AddTaskHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newTask models.Task

		// Get the user from the context
		idUser, ok := c.Locals("idUser").(float64)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("JWT not correct.")
		}

		if err := c.BodyParser(&newTask); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := validateTaskInput(&newTask); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		newTask.FkIdUser = int(idUser)
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
		var taskToUpdate models.Task
		var taskDb models.Task

		// Get the user from the context
		idUser, ok := c.Locals("idUser").(float64)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("JWT not correct.")
		}

		if err := c.BodyParser(&taskToUpdate); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := validateTaskInput(&taskToUpdate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		reqQuery := c.Queries()
		idTask := reqQuery["idTask"]

		if err := db.Where("tasks.id = ? and tasks.fk_id_user = ?", idTask, idUser).First(&taskDb).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Handle "record not found" error
				return c.Status(fiber.StatusNotFound).SendString("Task not found")
			} else {
				// Handle other errors
				return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving task from the database")
			}
		}

		if !taskToUpdate.CurrentTaskDate.IsZero() {
			taskDb.CurrentTaskDate = taskToUpdate.CurrentTaskDate
		}

		if taskToUpdate.Title != "" {
			taskDb.Title = taskToUpdate.Title
		}

		if taskToUpdate.Description != "" {
			taskDb.Description = taskToUpdate.Description
		}

		if taskToUpdate.FkIdCatCategory > 0 {
			taskDb.FkIdCatCategory = taskToUpdate.FkIdCatCategory
		}

		taskDb.FkIdUser = int(idUser)

		if taskToUpdate.FkIdCatStatus > 0 {
			taskDb.FkIdCatStatus = taskToUpdate.FkIdCatStatus
			taskDb.LastStatusChange = time.Now()
		}

		db.Save(taskDb)
		return c.JSON("updated")
	}
}

func ListTasksHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from the context
		idUser, ok := c.Locals("idUser").(float64)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("JWT not correct.")
		}

		var tasksDb []models.TaskApiResponse

		currentTime := time.Now()
		tomorrowTime := currentTime.Add(24 * time.Hour)

		err := db.Table("tasks").
			Select("tasks.id, tasks.title, tasks.description, tasks.fk_id_cat_status, tasks.current_task_date, cat_statuses.name as Name").
			Joins("left join users on users.id = tasks.fk_id_user").
			Joins("left join cat_statuses on cat_statuses.id = tasks.fk_id_cat_status").
			Where("users.id = ? and tasks.logical_delete = ?", idUser, false).
			Order("tasks.created_at desc").
			Scan(&tasksDb).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Handle "record not found" error
				return c.Status(fiber.StatusNotFound).SendString("Tasks not found")
			}else {
			  return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving users from the database")
      }
		}

		todayTasks := []models.TaskApiResponse{}
		tomorrowTasks := []models.TaskApiResponse{}

		for _, task := range tasksDb {
			currentDbTaskDate, err := time.Parse(time.RFC3339Nano, task.CurrentTaskDate)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, parsing dates.")
			}

      if currentDbTaskDate.Day() <= currentTime.Day() && currentDbTaskDate.Month() <= currentTime.Month() && currentDbTaskDate.Year() <= currentTime.Year() {
        if task.FkIdCatStatus == 1 {
				  todayTasks = append(todayTasks, task)
        }
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
