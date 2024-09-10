package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/types"
)

func (db *DB) GetCategories(ids []string, limit uint64, offset uint64, allow bool) (map[string]any, error) {

	if allow {
		query, args := db.QueryBuilder("SELECT * FROM categories c LEFT JOIN products p ON c.id = p.category_id WHERE c.id IN (%s)", ids, 0, 0)
		log.Info(query)
		rows, err := db.Conn.Query(query, args...)

		var results = make(map[string]any)
		if err != nil {
			return results, err
		}

		var categories = make(map[string]types.Category)
		var products = make(map[string][]types.Product)

		for rows.Next() {
			var prod types.Product
			var category types.Category
			err := rows.Scan(&category.Id, &category.Name, &prod.Id, &prod.Name, &prod.Stock, &prod.Price, &prod.Image, &prod.CategoryId, &prod.Created)
			if err != nil {
				return results, err
			}
			prodKey := fmt.Sprintf("%d", prod.CategoryId)
			catKey := fmt.Sprintf("%d", category.Id)
			products[prodKey] = append(products[prodKey], prod)
			categories[catKey] = category
		}
		if err = rows.Err(); err != nil {
			return results, err
		}

		for k, v := range products {
			results[k] = map[string]any{
				"category": categories[k],
				"products": v,
			}
			fmt.Println(results)
		}

		return results, nil
	}

	var categories = make(map[string]interface{})
	query, args := db.QueryBuilder("SELECT * FROM categories WHERE id IN (%s)", ids, 0, 0)
	log.Info(query)
	rows, err := db.Conn.Query(query, args...)

	if err != nil {
		return categories, err
	}

	defer rows.Close()

	for rows.Next() {
		var cat types.Category
		err := rows.Scan(&cat.Id, &cat.Name)
		if err != nil {
			return categories, err
		}
		catKey := fmt.Sprintf("%d", cat.Id)
		categories[catKey] = fiber.Map{"category": cat}
	}
	if err = rows.Err(); err != nil {
		return categories, err
	}
	return categories, nil
}

func (db *DB) GetAllCategories(limit uint64, offset uint64) ([]types.Category, error) {
	query, args := db.QueryBuilder("SELECT * FROM categories LIMIT $%d OFFSET $%d", nil, limit, offset)
	log.Info(query)
	rows, err := db.Conn.Query(query, args...)
	var categories []types.Category
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat types.Category
		err := rows.Scan(&cat.Id, &cat.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return categories, err
	}
	return categories, nil
}
