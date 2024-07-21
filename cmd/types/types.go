package types

import "time"

type Product struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Stock      int       `json:"stock"`
	Price      float32   `json:"price"`
	Image      string    `json:"image"`
	CategoryId int       `json:"category_id"`
	Created    time.Time `json:"created_at"`
	Updated    time.Time `json:"updated_at"`
}

type Order struct {
	Id         int       `json:"id"`
	CustomerId int       `json:"customer_id"`
	ProductId  int       `json:"product_id"`
	Created    time.Time `json:"created_at"`
	Updated    time.Time `json:"updated_at"`
}

type Customer struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
}

type Address struct {
	Id         int    `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	CustomerID uint   `json:"customer_id"`
}

type Category struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ProductId int    `json:"product_id"`
}
