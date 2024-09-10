package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
)

func RootHandler(c *fiber.Ctx) error {
	method := c.Method()
	log.Infof("Received a %s request for root", method)
	return c.SendString("Hello from root /")
}
