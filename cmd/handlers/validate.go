package handlers

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/types"
)

func Validate(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.User)
	if !ok {
		log.Info("User not found")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	fmt.Println("id", user.Id)
	return c.JSON(user)

}
