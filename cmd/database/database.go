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

func (db *DB) QueryBuilder(queryString string, ids []string, limit uint64, offset uint64) (string, []any) {

	var query string
	var args []interface{}
	if ids != nil {
		placeholders := make([]string, len(ids))
		for i := range ids {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		query = fmt.Sprintf(queryString, strings.Join(placeholders, ","))
		args = make([]interface{}, len(ids))
	} else {
		query = fmt.Sprintf(queryString, len(ids)+1, len(ids)+2)
		args = make([]interface{}, len(ids)+2)
	}

	for i, id := range ids {
		args[i] = id
	}

	if ids == nil {
		args[len(ids)] = limit
		args[len(ids)+1] = offset
	}
	fmt.Println("args: ", args)
	return query, args
}

func (db *DB) InsertOne(tableName string, item interface{}) (int, error) {
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
	log.Info(query)
	var pk int
	err := db.Conn.QueryRow(query, values...).Scan(&pk)
	return pk, err
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
	var pk int
	query += " RETURNING id"
	log.Info(query)
	err := db.Conn.QueryRow(query, values...).Scan(&pk)
	return pk, err
}

func (db *DB) UpdateOne(tableName string, id int, item interface{}) (int, error) {
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)
	addValues := []string{}

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			addValues = append(addValues, fmt.Sprintf("%s = %v", t.Field(i).Tag.Get("json"), v.Field(i).Interface()))
		}
	}

	var pk int
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d RETURNING id", tableName, strings.Join(addValues, ","), id)
	log.Info(query)
	err := db.Conn.QueryRow(query).Scan(&pk)
	return pk, err
}

func (db *DB) DeleteOne(tableName string, id string) (int, error) {
	var pk int
	err := db.Conn.QueryRow("DELETE FROM $1 where id = $2 RETURNING id", tableName, id).Scan(&pk)
	return pk, err
}

func (db *DB) DeleteMany(tableName string, where string, ids []string) ([]int, error) {
	if ids == nil {
		return nil, nil
	}

	queryString := fmt.Sprintf("DELETE FROM %s where %s", tableName, where)
	placeholders := make([]string, len(ids))

	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("%s IN (%s) RETURNING id", queryString, strings.Join(placeholders, ","))
	log.Info(query)
	idA := make([]interface{}, len(ids))
	for i, v := range ids {
		idA[i] = v
	}

	rows, err := db.Conn.Query(query, idA...)

	if err != nil {
		return nil, err
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
		price NUMERIC(9,2) NOT NULL,
		image VARCHAR(255) NOT NULL,
		category_id INT NOT NULL,
		created_at timestamp DEFAULT NOW(),
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
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
		created_at TIMESTAMPTZ DEFAULT NOW()
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
	created_at TIMESTAMPTZ DEFAULT NOW(),
	FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
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
		FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateCategories(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS  categories(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUsers(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		password TEXT NOT NULL,
		role ROLE NOT NULL,
		UNIQUE(username)
		)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	typeQuery := `CREATE TYPE IF NOT EXISTS role AS ENUM('admin', 'manager', 'sales');`
	_, err = db.Exec(typeQuery)
	if err != nil {
		log.Warn(err)
	}
}
