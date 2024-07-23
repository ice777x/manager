package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/database"
)

func DbWare(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	}
}
