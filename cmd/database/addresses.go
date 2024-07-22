package database

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ice777x/pmanager/cmd/types"
	"strings"
)

func (db *DB) GetAddresses(addressIds []string, limit uint64, offset uint64) ([]types.Address, error) {
	query := fmt.Sprintf("SELECT * FROM addresses WHERE id IN (%s) LIMIT %d OFFSET %d", strings.Join(addressIds, ","), limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var addresses []types.Address
	for rows.Next() {
		var address types.Address
		err := rows.Scan(&address.Id, &address.Street, &address.City, &address.State, &address.ZipCode, &address.CustomerID)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, address)
	}
	if err = rows.Err(); err != nil {
		return addresses, err
	}
	return addresses, nil
}

func (db *DB) GetAllAddresses(limit uint64, offset uint64) ([]types.Address, error) {
	query := fmt.Sprintf("SELECT * FROM addresses LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var addresses []types.Address
	for rows.Next() {
		var address types.Address
		err := rows.Scan(&address.Id, &address.Street, &address.City, &address.State, &address.ZipCode, &address.CustomerID)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, address)
	}
	if err = rows.Err(); err != nil {
		return addresses, err
	}
	return addresses, nil
}
