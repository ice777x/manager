package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ice777x/manager/cmd/types"
	"strings"
)

func (db *DB) GetOrders(orderIds []string, limit uint64, skip uint64) ([]types.Order, error) {
	query := fmt.Sprintf("SELECT * FROM orders WHERE id IN (%s) LIMIT %d OFFSET %d", strings.Join(orderIds, ","), limit, skip)

	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err := rows.Scan(&order.Id, &order.ProductId, &order.CustomerId, &order.Created, &order.Updated)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return orders, err
	}
	return orders, nil
}

func (db *DB) GetAllOrders(limit uint64, offset uint64) ([]types.Order, error) {
	query := fmt.Sprintf("SELECT * FROM orders LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err := rows.Scan(&order.Id, &order.ProductId, &order.CustomerId, &order.Created, &order.Updated)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)

		if err = rows.Err(); err != nil {
			return orders, err
		}
	}
	return orders, nil
}
