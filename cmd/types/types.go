package types

import "time"

type Product struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Stock      int       `json:"stock"`
	Price      float32   `json:"price"`
	Image      string    `json:"image"`
	CategoryId uint      `json:"category_id"`
	Created    time.Time `json:"created_at"`
}

type Order struct {
	Id         int       `json:"id"`
	CustomerId int       `json:"customer_id"`
	ProductId  int       `json:"product_id"`
	Created    time.Time `json:"created_at"`
}

type Customer struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Created   time.Time `json:"created_at"`
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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CustomerAddress struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Address   Address   `json:"address"`
	Created   time.Time `json:"created_at"`
}

type OrderProdCust struct {
	ID       int       `json:"id"`
	Product  Product   `json:"product"`
	Category Category  `json:"category"`
	Customer Customer  `json:"customer"`
	Address  Address   `json:"address"`
	Created  time.Time `json:"created_at"`
}
