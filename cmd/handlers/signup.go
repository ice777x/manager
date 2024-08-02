package handlers

import (
	"slices"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
	"golang.org/x/crypto/bcrypt"
)

var Role = []string{"admin", "manager", "sales", "customer"}

func SignUp(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)
	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var user types.User
	if err := c.BodyParser(&user); err != nil {
		log.Errorf("Signup Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid parameters",
		})
	}
	if user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse username",
		})
	}
	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse password",
		})
	}

	if user.Role == "" || !slices.Contains(Role, user.Role) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse Role",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Errorf("Signup Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to generate password",
		})
	}
	user.Password = string(hash)
	_, err = db.InsertOne("users", user)
	if err != nil {
		log.Errorf("Signup Handler: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusOK)
}
