package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
)

func ProductItem(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	qs := c.Queries()

	var limit uint64 = 10
	var skip uint64 = 0
	var err error
	id := qs["id"]

	if limitStr := qs["limit"]; limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)

		if err != nil {
			log.Errorf("Product Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid limit parameter",
			})
		}
	}

	if skipStr := qs["skip"]; skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)

		if err != nil {
			log.Errorf("Product Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	if id == "" {
		products, err := db.GetAllProduct(limit, skip)
		if err != nil {
			log.Errorf("Product Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Failed to retireve products",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"results": products,
		})
	}

	products, err := db.GetProduct(strings.Split(strings.Trim(id, ","), ","), limit, skip)
	if err != nil {
		log.Errorf("Product Handler: %v", err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "No data for query",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"results": products,
	})
}

func ProductCreate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req []types.Product
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Product Create Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	products := make([]interface{}, len(req))
	for i, product := range req {
		product.Created = time.Now()
		products[i] = product
	}

	_, err := db.InsertMany("products", products)
	if err != nil {
		log.Errorf("Product Create Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to insert products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 201,
		"result": "Products added successfully!",
	})

}

func ProductUpdate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Product Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse id.",
		})
	}

	var req types.Product
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Product Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	pk, err := db.UpdateOne("products", id, req)

	if err != nil {
		log.Errorf("Product Update Handler: %v", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to parse request body as JSON. Please check your input and try again.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": fmt.Sprintf("ID=%d is successfully updated!", pk),
	})
}

func ProductDelete(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Query("id")
	id := strings.Split(strings.Trim(idStr, ","), ",")

	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "id parameter not found",
		})
	}

	pk, err := db.DeleteMany("products", "id", id)
	if err != nil {
		log.Errorf("Product Delete Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to delete products",
		})
	}
	log.Infof("Delete item from product %s", id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": pk,
	})
}
