package handlers

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/database"
	"github.com/ice777x/manager/cmd/types"
)

func CategoryItem(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	qs := c.Queries()

	var limit uint64 = 10
	var skip uint64 = 0
	var allowProduct = false
	var err error

	id := qs["id"]

	if limitStr := qs["limit"]; limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			log.Errorf("Category Handler %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid limit parameter",
			})
		}
	}

	if skipStr := qs["skip"]; skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			log.Errorf("Category Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	if id == "" {
		categories, err := db.GetAllCategories(limit, skip)
		if err != nil {
			log.Errorf("Category Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Failed to retrieve categories",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": 200,
			"result": categories,
		})
	}

	if allowProductStr := qs["product"]; allowProductStr != "" {
		allowProduct, err = strconv.ParseBool(allowProductStr)

		if err != nil {
			log.Errorf("Category Handler: %v", err)
			return c.JSON(fiber.Map{
				"status":  400,
				"message": "Invalid skip parameter",
			})
		}
	}

	categories, err := db.GetCategories(strings.Split(strings.Trim(id, ","), ","), limit, skip, allowProduct)
	if err != nil {
		log.Errorf("Category Handler: %v", err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "No data for query",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"results": categories,
	})
}

func CategoryCreate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	var req []types.Category
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Category Create Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	categories := make([]interface{}, len(req))
	for i, category := range req {
		categories[i] = category
	}

	pk, err := db.InsertMany("categories", categories)
	if err != nil {
		log.Errorf("Category Create Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to insert categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": pk,
	})

}

func CategoryUpdate(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Category Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse id.",
		})
	}
	var req types.Category
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Category Update Handler: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to parse request body as JSON. Please check your data and try again.",
		})
	}

	pk, err := db.UpdateOne("categories", id, req)
	if err != nil {
		log.Errorf("Category Update Handler: %v", err)
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

func CategoryDelete(c *fiber.Ctx) error {

	db, ok := c.Locals("db").(*database.DB)

	if !ok {
		log.Fatal("Problem in database connection!")
	}

	idStr := c.Params("id")
	id := strings.Split(strings.Trim(idStr, ","), ",")
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Id parameter not found",
		})
	}

	log.Infof("DELETE ITEM FROM %s", id)

	res, err := db.DeleteMany("categories", "id", id)
	if err != nil {
		log.Errorf("Category Delete Handler: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Failed to delete categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 200,
		"result": res,
	})
}
