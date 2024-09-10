package database

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/ice777x/manager/cmd/types"
)

func (db *DB) GetProduct(ids []string, limit uint64, offset uint64) ([]interface{}, error) {
	query, args := db.QueryBuilder("SELECT * FROM products p INNER JOIN categories c ON c.id = p.category_id WHERE p.id IN (%s)", ids, 0, 0)

	log.Info(query)
	rows, err := db.Conn.Query(query, args...)
	var products []any
	if err != nil {
		return products, errors.New("id is invalid")
	}
	defer rows.Close()
	for rows.Next() {
		var prod types.Product
		var cat types.Category
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Stock, &prod.Price, &prod.Image, &prod.CategoryId, &prod.Created, &cat.Id, &cat.Name)
		if err != nil {
			log.Debug(err)
			return products, err
		}
		products = append(products, fiber.Map{
			"product":  prod,
			"category": cat,
		})
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}

func (db *DB) GetAllProduct(limit uint64, offset uint64) ([]interface{}, error) {
	query, args := db.QueryBuilder("SELECT * FROM products p INNER JOIN categories c ON c.id = p.category_id LIMIT $%d OFFSET $%d", nil, limit, offset)
	log.Info(query)
	rows, err := db.Conn.Query(query, args...)

	var products []any
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var prod types.Product
		var cat types.Category
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Stock, &prod.Price, &prod.Image, &prod.CategoryId, &prod.Created, &cat.Id, &cat.Name)
		if err != nil {
			log.Debug(err)
			return products, err
		}
		products = append(products, fiber.Map{
			"product":  prod,
			"category": cat,
		})
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}
