package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/simoncra/godoit/config"
	"github.com/simoncra/godoit/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// getUsersHandler retrieves a lsit of users
func GetUsersHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving users from the database")
		}

		return c.JSON(users)
	}
}

func CreateUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newUser models.User

		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := validateUserInput(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Get the user from the context
		userJwt := c.Locals("user").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		fmt.Println("Welcome 👋" + email)

		hashedPsswd, err := hashPassword(newUser.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
		}

		newUser.Password = hashedPsswd

		if err := db.Create(&newUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating a new user")
		}

		return c.JSON(newUser)
	}
}

func SignupHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newUser models.User

		// load configuration
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := validateUserInput(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		hashedPsswd, err := hashPassword(newUser.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
		}

		newUser.Password = hashedPsswd

		if err := db.Create(&newUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating a new user")
		}

		// Create the JWT claims, which includes the user ID and expiry time
		claims := jwt.MapClaims{
			"ID":    newUser.ID,
			"email": newUser.Email,
			"exp":   time.Now().Add(time.Hour * 24 * 1).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		tokenEncoded, err := token.SignedString([]byte(config.AppSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return the token
		return c.JSON(models.SignupResponse{
			Token: tokenEncoded,
		})
	}
}

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginRequest models.LoginRequest
		var user models.User

		// load configuration
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		if err := c.BodyParser(&loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				// User not found
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
			}

			// Database error
			return c.Status(fiber.StatusInternalServerError).SendString("Error querying the database")
		}

		if !validatePassword(loginRequest.Password, user.Password) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password.")
		}

		day := time.Hour * 24

		// Create the JWT claims, which includes the user ID and expiry time
		claims := jwt.MapClaims{
			"ID":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(day * 1).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		tokenEncoded, err := token.SignedString([]byte(config.AppSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return the token
		return c.JSON(models.LoginResponse{
			Token: tokenEncoded,
		})
	}
}

func validateUserInput(user *models.User) error {
	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		var validationErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation error: %s", validationErrors)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validatePassword(inputPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}
