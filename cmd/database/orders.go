package database

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/ice777x/manager/cmd/types"
)

func (db *DB) GetOrders(ids []string, limit uint64, skip uint64) ([]types.OrderProdCust, error) {
	query, args := db.QueryBuilder("SELECT o.id, p.*, cat.*, c.*, adr.*, o.created_at FROM orders o INNER JOIN products p ON p.id = o.product_id INNER JOIN customers c ON c.id = o.customer_id INNER JOIN categories cat ON cat.id = p.category_id INNER JOIN addresses adr USING(customer_id) WHERE o.id IN (%s)", ids, 0, 0)

	log.Info(query)

	rows, err := db.Conn.Query(query, args...)
	var orders []types.OrderProdCust
	if err != nil {
		fmt.Println(err)
		return orders, err
	}

	defer rows.Close()
	for rows.Next() {
		var order types.OrderProdCust
		err := rows.Scan(&order.ID, &order.Product.Id, &order.Product.Name, &order.Product.Stock, &order.Product.Price, &order.Product.Image, &order.Product.CategoryId, &order.Product.Created, &order.Category.Id, &order.Category.Name, &order.Customer.Id, &order.Customer.FirstName, &order.Customer.LastName, &order.Customer.Created, &order.Address.Id, &order.Address.Street, &order.Address.City, &order.Address.State, &order.Address.ZipCode, &order.Address.CustomerID, &order.Created)
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

func (db *DB) GetAllOrders(limit uint64, skip uint64) ([]types.OrderProdCust, error) {
	query, args := db.QueryBuilder("SELECT o.id, p.*, cat.*, c.*, adr.*,o.created_at FROM orders o INNER JOIN products p ON p.id = o.product_id INNER JOIN customers c ON c.id = o.customer_id INNER JOIN categories cat ON cat.id = p.category_id INNER JOIN addresses adr USING(customer_id) LIMIT $%d OFFSET $%d", nil, limit, skip)
	log.Info(query)
	rows, err := db.Conn.Query(query, args...)

	var orders []types.OrderProdCust
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order types.OrderProdCust
		err := rows.Scan(&order.ID, &order.Product.Id, &order.Product.Name, &order.Product.Stock, &order.Product.Price, &order.Product.Image, &order.Product.CategoryId, &order.Product.Created, &order.Category.Id, &order.Category.Name, &order.Customer.Id, &order.Customer.FirstName, &order.Customer.LastName, &order.Customer.Created, &order.Address.Id, &order.Address.Street, &order.Address.City, &order.Address.State, &order.Address.ZipCode, &order.Address.CustomerID, &order.Created)
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
