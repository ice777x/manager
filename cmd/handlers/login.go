package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)
	if !ok {
		log.Fatal("Problem in database connection!")
	}
	var body types.User
	if err := c.BodyParser(&body); err != nil {
		log.Errorf("Login Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}
	var user types.User
	query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", body.Username)

	err := db.Conn.QueryRow(query).Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		log.Errorf("Login Handler: %v", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Errorf("Login Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid email or password",
		})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Errorf("Login Handler: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.MaxAge = 3600 * 24 * 30
	cookie.Secure = false
	cookie.HTTPOnly = true
	c.Cookie(cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": tokenString,
	})
}
