package main

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/middleware"
	"github.com/ice777x/manager/cmd/utils"
	_ "github.com/lib/pq"
)

func initialize(db *sql.DB) {
	database.CreateCustomer(db)
	database.CreateAddress(db)
	database.CreateProducts(db)
	database.CreateCategories(db)
	database.CreateOrders(db)
	database.CreateUsers(db)
}

func main() {
	fmt.Println("MANAGER v0.1")
	utils.Config()

	con, err := utils.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db := &database.DB{Conn: con}
	initialize(db.Conn)

	app := fiber.New()
	app.Use(middleware.DbWare(db))
	app.Use("/api", middleware.AuthGuard)
	cmd.Router(app)
	defer db.Conn.Close()
}
