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

func CustomerItem(c *fiber.Ctx) error {
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
			log.Errorf("Category Delete Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid limit parameter",
			})
		}
	}

	if skipStr := qs["skip"]; skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)

		if err != nil {

			log.Errorf("Customer Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	if id == "" {
		customers, err := db.GetAllCustomers(limit, skip)
		if err != nil {
			log.Errorf("Customer Handler: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  500,
				"message": "Failed to retrieve customers",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": 200,
			"result": customers,
		})
	}
	customers, err := db.GetCustomers(strings.Split(strings.Trim(id, ","), ","))
	if err != nil {
		log.Errorf("Customer Handler: %v", err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "No data for query",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": customers,
	})
}

func CustomerCreate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req []types.Customer
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Customer Create Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	customers := make([]interface{}, len(req))
	for i, customer := range req {
		customer.Created = time.Now()
		customers[i] = customer
	}

	i, err := db.InsertMany("customers", customers)
	if err != nil {
		log.Errorf("Customer Create Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to insert customers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": i,
	})

}

func CustomerUpdate(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Customer Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse id.",
		})
	}

	var req types.Customer
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Customer Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}
	pk, err := db.UpdateOne("customers", id, req)
	if err != nil {
		log.Errorf("Customer Update Handler: %v", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to update customer.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": pk,
	})
}

func CustomerDelete(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Query("id")
	id := strings.Split(strings.Trim(idStr, ","), ",")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "id parameter not found",
		})
	}
	fmt.Println(id)
	res, err := db.DeleteMany("customers", "id", id)
	if err != nil {
		log.Errorf("Customer Delete Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to delete customers",
		})
	}
	log.Infof("Delete item from customers %d", res)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": res,
	})
}
