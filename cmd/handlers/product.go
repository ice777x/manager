package handlers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/pmanager/cmd/database"
	"github.com/ice777x/pmanager/cmd/types"
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

			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid limit parameter",
			})
		}
	}

	if skipStr := qs["skip"]; skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)

		if err != nil {
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	if id == "" {
		products, err := db.GetAllProduct(limit, skip)
		if err != nil {
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Failed to retireve products",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": 200,
			"result": products,
		})
	}

	products, err := db.GetProduct(strings.Split(strings.Trim(id, ","), ","), limit, skip)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "No data for query",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": products,
	})
}

func ProductCreate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req []types.Product
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}

	products := make([]interface{}, len(req))
	for i, product := range req {
		product.Created = time.Now()
		product.Updated = time.Now()
		products[i] = product
	}

	i, err := db.InsertMany("products", products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to insert products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": i,
	})

}

func ProductUpdate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req types.Product
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": err.Error(),
		})
	}

	pk, err := db.UpdateOne("products", req)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": pk,
	})
}

func ProductDelete(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("ids")
	id := strings.Split(strings.Trim(idStr, ","), ",")

	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "id parameter not found",
		})
	}

	log.Infof("DELETE ITEM FROM %s", id)

	res, err := db.DeleteMany("products", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to delete products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": res,
	})
}
