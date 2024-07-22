package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ice777x/pmanager/cmd/types"
)

func (db *DB) GetProduct(productIds []string, limit uint64, offset uint64) ([]types.Product, error) {
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s) LIMIT %d OFFSET %d", strings.Join(productIds, ","), limit, offset)
	rows, err := db.Conn.Query(query)
	var products []types.Product
	if err != nil {
		return products, errors.New("id is invalid")
	}
	defer rows.Close()
	for rows.Next() {
		var prod types.Product
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Stock, &prod.Price, &prod.Image, &prod.CategoryId, &prod.Created, &prod.Updated)
		if err != nil {
			return products, err
		}
		products = append(products, prod)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}

func (db *DB) GetAllProduct(limit uint64, offset uint64) ([]types.Product, error) {
	query := fmt.Sprintf("SELECT * FROM products LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Conn.Query(query)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var prod types.Product
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Stock, &prod.Price, &prod.Image, &prod.CategoryId, &prod.Created, &prod.Updated)
		if err != nil {
			return products, err
		}
		products = append(products, prod)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}
	return products, nil
}
