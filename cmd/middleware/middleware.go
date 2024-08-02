package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
)

func DbWare(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	}
}

func AuthGuard(c *fiber.Ctx) error {
	if c.OriginalURL() == "/api/signup" || c.OriginalURL() == "/api/login" {
		return c.Next()
	}

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	authToken := c.Cookies("Authorization")
	if authToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("xd")
			c.ClearCookie("Authorization")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", claims["sub"])
		var user types.User
		err := db.Conn.QueryRow(query).Scan(&user.Id, &user.Username, &user.Password, &user.Role)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  400,
				"message": "User fot found",
			})
		}
		c.Locals("user", user)
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
