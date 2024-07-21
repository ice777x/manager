package main

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/pmanager/cmd/database"
	"github.com/ice777x/pmanager/cmd/handlers"
	"github.com/ice777x/pmanager/cmd/middleware"
	"github.com/ice777x/pmanager/cmd/utils"
	_ "github.com/lib/pq"
)

func initialize(db *sql.DB) {
	database.CreateCustomer(db)
	database.CreateAddress(db)
	database.CreateProducts(db)
	database.CreateCategories(db)
	database.CreateOrders(db)
}

func main() {
	fmt.Println("HELLO FROM PRODUCT_MANAGER_APP")
	utils.Config()

	con, err := utils.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	db := &database.DB{Conn: con}
	initialize(db.Conn)

	app := fiber.New()
	app.Use(middleware.DbWare(db))

	Router(app)
	defer db.Conn.Close()
}

func Router(app *fiber.App) {
	app.Route("/api", func(api fiber.Router) {
		api.Get("/", handlers.RootHandler)

		api.Get("/product", handlers.ProductItem)
		api.Post("/product", handlers.ProductCreate)
		api.Put("/product", handlers.ProductUpdate)
		api.Delete("/product", handlers.ProductDelete)

		api.Get("/order", handlers.OrderItem)
		api.Post("/order", handlers.OrderCreate)
		api.Put("/order", handlers.OrderUpdate)
		api.Delete("/order", handlers.OrderDelete)

		api.Get("/customer", handlers.CustomerItem)
		api.Post("/customer", handlers.CustomerCreate)
		api.Put("/customer", handlers.CustomerUpdate)
		api.Delete("/customer", handlers.CustomerDelete)

		api.Get("/address", handlers.AddressItem)
		api.Post("/address", handlers.AddressCreate)
		api.Put("/address", handlers.AddressUpdate)
		api.Delete("/address", handlers.AddressDelete)

		api.Get("/category", handlers.CategoryItem)
		api.Post("/category", handlers.CategoryCreate)
		api.Put("/category", handlers.CategoryUpdate)
		api.Delete("/category", handlers.CategoryDelete)
	}, "api")

	log.Info("Starting server on port 3000")
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
