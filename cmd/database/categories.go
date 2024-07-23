package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ice777x/manager/cmd/types"
	"strings"
)

func (db *DB) GetCategories(categoryIds []string, limit uint64, offset uint64) ([]types.Category, error) {
	query := fmt.Sprintf("SELECT * FROM categories WHERE id IN (%s) LIMIT %d OFFSET %d", strings.Join(categoryIds, ","), limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var categories []types.Category
	for rows.Next() {
		var cat types.Category
		err := rows.Scan(&cat.Id, &cat.Name, &cat.ProductId)
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

func (db *DB) GetAllCategories(limit uint64, offset uint64) ([]types.Category, error) {
	query := fmt.Sprintf("SELECT * FROM categories LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []types.Category
	for rows.Next() {
		var cat types.Category
		err := rows.Scan(&cat.Id, &cat.Name, &cat.ProductId)
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
