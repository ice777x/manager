package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ice777x/manager/cmd/types"
)

func (db *DB) GetCustomers(ids []string) ([]types.CustomerAddress, error) {
	query, args := db.QueryBuilder("SELECT * FROM customers c LEFT JOIN addresses adr ON c.id = adr.customer_id WHERE c.id IN (%s)", ids, 0, 0)
	fmt.Println(args)
	log.Info(query)

	rows, err := db.Conn.Query(query, args...)
	fmt.Println(err)
	var customers []types.CustomerAddress

	if err != nil {
		return customers, err
	}
	defer rows.Close()

	for rows.Next() {
		var cust types.CustomerAddress
		err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Created, &cust.Address.Id, &cust.Address.Street, &cust.Address.City, &cust.Address.State, &cust.Address.ZipCode, &cust.Address.CustomerID)
		if err != nil {
			return customers, err
		}
		customers = append(customers, cust)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil
}

func (db *DB) GetAllCustomers(limit uint64, offset uint64) ([]types.CustomerAddress, error) {
	query, args := db.QueryBuilder("SELECT * FROM customers c LEFT JOIN addresses adr ON c.id = adr.customer_id LIMIT $%d OFFSET $%d", nil, limit, offset)

	log.Info(query)
	rows, err := db.Conn.Query(query, args...)
	var customers []types.CustomerAddress

	if err != nil {
		return customers, err
	}
	defer rows.Close()

	for rows.Next() {
		var cust types.CustomerAddress
		err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Created, &cust.Address.Id, &cust.Address.Street, &cust.Address.City, &cust.Address.State, &cust.Address.ZipCode, &cust.Address.CustomerID)
		if err != nil {
			return customers, err
		}
		customers = append(customers, cust)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil
}
