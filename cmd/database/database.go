package database

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	"reflect"
	"strings"
)

type DB struct {
	Conn *sql.DB
}

func (db *DB) InsertOne(tableName string, item interface{}) int {
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	fieldNames := make([]string, v.NumField()-1)
	placeholders := make([]string, v.NumField()-1)
	values := make([]interface{}, v.NumField()-1)

	for i := 1; i < v.NumField(); i++ {
		fieldNames[i-1] = t.Field(i).Tag.Get("json")
		placeholders[i-1] = fmt.Sprintf("$%d", i)
		values[i-1] = v.Field(i).Interface()
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)
	var pk int
	err := db.Conn.QueryRow(query, values...).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}

func (db *DB) InsertMany(tableName string, items []interface{}) (int, error) {
	if len(items) == 0 {
		return 0, nil
	}
	val := reflect.ValueOf(items[0])
	tag := reflect.TypeOf(items[0])
	fieldNames := make([]string, val.NumField()-1)
	placeholders := make([]string, len(items))
	values := make([]interface{}, (val.NumField()-1)*len(items))

	for i := 1; i < val.NumField(); i++ {
		fieldNames[i-1] = tag.Field(i).Tag.Get("json")
	}

	for x, item := range items {
		v := reflect.ValueOf(item)

		for i := 1; i < v.NumField(); i++ {
			placeholders[x] += fmt.Sprintf("$%d", (i + (x)*(v.NumField()-1)))
			if i != v.NumField()-1 {
				placeholders[x] += ","
			}
			values[(i+(x)*(v.NumField()-1))-1] = v.Field(i).Interface()
		}
	}
	query := fmt.Sprintf("INSERT INTO %s (%s)  VALUES ", tableName, strings.Join(fieldNames, ","))
	for i, row := range placeholders {
		query += "(" + row + ")"
		if i != len(placeholders)-1 {
			query += ","
		}
	}
	result, err := db.Conn.Exec(query, values...)
	if err != nil {
		log.Fatal(err)
	}

	lId, _ := result.LastInsertId()
	return int(lId), err
}

func (db *DB) UpdateOne(tableName string, item interface{}) (int, error) {
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)
	addValues := []string{}

	var idStr any
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			if t.Field(i).Tag.Get("json") == "id" {
				idStr = v.Field(i).Interface()
				continue
			}
			addValues = append(addValues, fmt.Sprintf("%s = %v", t.Field(i).Tag.Get("json"), v.Field(i).Interface()))
		}
	}
	var pk int
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %v RETURNING id", tableName, strings.Join(addValues, ","), idStr)
	fmt.Println(query)
	err := db.Conn.QueryRow(query).Scan(&pk)
	return pk, err
}

func (db *DB) DeleteOne(tableName string, id string) (int, error) {
	query := fmt.Sprintf("DELETE FROM %s where id = %s RETURNING id", tableName, id)
	var dId int
	err := db.Conn.QueryRow(query).Scan(&dId)
	return dId, err
}

func (db *DB) DeleteMany(tableName string, ids []string) ([]int, error) {
	query := fmt.Sprintf("DELETE FROM %s where id IN (%s) RETURNING id", tableName, strings.Join(ids, ","))
	rows, err := db.Conn.Query(query)
	if err != nil {
		return []int{}, err
	}
	var pks []int
	defer rows.Close()
	for rows.Next() {
		var pk int
		err := rows.Scan(&pk)
		if err != nil {
			return pks, err
		}
		pks = append(pks, pk)
	}
	return pks, err
}

func CreateProducts(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		stock INT NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		image VARCHAR(255) NOT NULL,
		category_id NUMERIC NOT NULL,
		updated_at timestamp DEFAULT NOW(),
		created_at timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateCustomer(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateOrders(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	customer_id INT NOT NULL,
	product_id INT NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	FOREIGN KEY (customer_id) REFERENCES customers(id),
	FOREIGN KEY (product_id) REFERENCES products(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateAddress(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS addresses (
		id SERIAL PRIMARY KEY,
		street VARCHAR(255) NOT NULL,
		city VARCHAR(255) NOT NULL,
		state VARCHAR(255) NOT NULL,
		zip_code VARCHAR(255) NOT NULL,
		customer_id INT NOT NULL,
		FOREIGN KEY (customer_id) REFERENCES customers(id) 
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateCategories(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS  categories(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		product_id INT NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products(id) 
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
