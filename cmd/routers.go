package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/handlers"
)

func Router(app *fiber.App) {
	app.Route("/api", func(api fiber.Router) {
		api.Get("/", handlers.RootHandler)
		api.Post("/signup", handlers.SignUp).Name("Signup")
		api.Post("/login", handlers.Login).Name("Login")
		api.Get("/validate", handlers.Validate)

		api.Get("/product", handlers.ProductItem)
		api.Post("/product", handlers.ProductCreate)
		api.Put("/product/:id", handlers.ProductUpdate)
		api.Delete("/product", handlers.ProductDelete)

		api.Get("/order", handlers.OrderItem)
		api.Post("/order", handlers.OrderCreate)
		api.Put("/order/:id", handlers.OrderUpdate)
		api.Delete("/order", handlers.OrderDelete)

		api.Get("/customer", handlers.CustomerItem)
		api.Post("/customer", handlers.CustomerCreate)
		api.Put("/customer/:id", handlers.CustomerUpdate)
		api.Delete("/customer", handlers.CustomerDelete)

		api.Get("/category", handlers.CategoryItem)
		api.Post("/category", handlers.CategoryCreate)
		api.Put("/category/:id", handlers.CategoryUpdate)
		api.Delete("/category", handlers.CategoryDelete)
	}, "api")

	log.Info("Starting server on port 3000")
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
