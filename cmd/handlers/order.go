package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
)

func OrderItem(c *fiber.Ctx) error {
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

			log.Errorf("Order Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid limit parameter",
			})
		}
	}

	if skipStr := qs["skip"]; skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)

		if err != nil {

			log.Errorf("Order Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	if id == "" {
		orders, err := db.GetAllOrders(limit, skip)
		if err != nil {
			log.Errorf("Order Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Failed to get orders",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": 200,
			"result": orders,
		})
	}

	orders, err := db.GetOrders(strings.Split(strings.Trim(id, ","), ","), limit, skip)
	if err != nil {
		log.Errorf("Order Handler: %v", err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "No data for query",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": orders,
	})
}

func OrderCreate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req []types.Order
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Order Create Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	orders := make([]interface{}, len(req))
	for i, order := range req {
		order.Created = time.Now()
		orders[i] = order
	}

	pk, err := db.InsertMany("orders", orders)
	if err != nil {
		log.Errorf("Order Create Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to insert orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 201,
		"result": pk,
	})

}

func OrderUpdate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Order Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse id.",
		})
	}

	var req types.Order
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Order Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your input and try again.",
		})
	}

	_, err = db.UpdateOne("orders", id, req)
	if err != nil {
		log.Errorf("Order Update Handler: %v", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": "Orders added successfully!",
	})
}

func OrderDelete(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("id")
	id := strings.Split(strings.Trim(idStr, ","), ",")
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid id parameter",
		})
	}

	log.Infof("DELETE ITEM FROM %s", id)

	res, err := db.DeleteMany("orders", "id", id)

	if err != nil {
		log.Errorf("Order Delete  Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to delete orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": res,
	})
}
