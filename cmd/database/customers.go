package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ice777x/pmanager/cmd/types"
	"strings"
)

func (db *DB) GetCustomers(customerIds []string, limit uint64, offset uint64) ([]types.Customer, error) {
	query := fmt.Sprintf("SELECT * FROM customers WHERE id IN (%s) LIMIT %d OFFSET %d", strings.Join(customerIds, ","), limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var customers []types.Customer
	for rows.Next() {
		var cust types.Customer
		err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Created, &cust.Updated)
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

func (db *DB) GetAllCustomers(limit uint64, offset uint64) ([]types.Customer, error) {
	query := fmt.Sprintf("SELECT * FROM customers LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var customers []types.Customer
	for rows.Next() {
		var cust types.Customer
		err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Created)
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
