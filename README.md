<p align="center">
  <img src="https://cdn-icons-png.flaticon.com/512/6295/6295417.png" width="100" />
</p>
<p align="center">
    <h1 align="center">MANAGER</h1>
</p>
<p align="center">
    <em>Manager is a product and order management API coded in Go. This API provides a simple and efficient way to manage products and orders in a system.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/ice777x/manager?style=flat&color=4011ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/ice777x/manager?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/ice777x/manager?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/ice777x/manager?style=flat&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white" alt="PostgreSQL" />
</p>
<hr>

## 📦 Features

- Product Management: Add, update, delete, and retrieve product information.
- Order Management: Create, update, delete, and retrieve orders.

## ⚙️ Installation

To install and run the project, follow these steps:

1. **Clone the repository:**
    ```bash
    git clone https://github.com/ice777x/manager.git
    cd manager
    ```
    
2. **Set up your environment variables:**
    Create a `.env` file in the root directory and add the following:
    ```plaintext
    DB_NAME=db_name
    DB_USER=user
    DB_PASS=passs
    DB_HOST=host
    ```
    
3. **Build the project:**
    ```bash
    go build
    ```

4. **Run the project:**
    ```bash
    ./manager
    ```

## 🚀 Usage

Once the project is running, you can use the following endpoints to interact with the API:
### Users
- **POST** `/api/signup`: Create a new user.
- **POTS** `/api/login`: Login as user.

### Products

- **GET** `/api/products`: Retrieve all products.
- **GET** `/api/products?id=1,2,3`: Retrieve a product by ID.
- **POST** `/api/products`: Create a new product.
- **PUT** `/api/products/1`: Update a product by ID.
- **DELETE** `/api/products?id=5,6`: Delete a product by ID.

### Orders

- **GET** `/api/orders`: Retrieve all orders.
- **GET** `/api/orders?id=1,2,3,4`: Retrieve an order by ID.
- **POST** `/api/orders`: Create a new order.
- **PUT** `/api/orders/1`: Update an order by ID.
- **DELETE** `/api/orders?id=4`: Delete an order by ID.

### Categories

- **GET** `/api/categories`: Retrieve all orders.
- **GET** `/api/categories?id=1,2,3,4`: Retrieve an order by ID.
- **POST** `/api/categories`: Create a new order.
- **PUT** `/api/categories/1`: Update an order by ID.
- **DELETE** `/api/categories?id=4`: Delete an order by ID.

### Customers

- **GET** `/api/customers`: Retrieve all orders.
- **GET** `/api/customers?id=1,2,3,4`: Retrieve an order by ID.
- **POST** `/api/customers`: Create a new order.
- **PUT** `/api/customers/1`: Update an order by ID.
- **DELETE** `/api/customers?id=4`: Delete an order by ID.

### Example Requests

- **Get a specific product:**
    ```bash
    curl http://localhost:3000/products?id=1
    ```

    Example Response:
    ```json
    {
      "result": [
      {
        "category": {
                "id": 1,
                "name": "Masonry & Precast"
            },
            "product": {
                "id": 1,
                "name": "Crawler",
                "stock": 3273,
                "price": 28496.33,
                "image": "http://dummyimage.com/152x100.png/5fa2dd/ffffff",
                "category_id": 1,
                "created_at": "2024-07-30T10:49:37.984252Z"
            }
      }
    ],
      "status": 200
    }
    ```
    
- **Create a new product:**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '[{"name":"Çivi", "stock":99, "price":0.25, "image":"civi.jpg", "category_id":1}]' http://localhost:3000/products
    ```
    ```json
    {
      "result": 2,
      "status": 200
    }
    ```
    
- **Update a product:**
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Product Name", "price":150}' http://localhost:3000/products/8
    {
      "result": 8,
      "status": 200
    }
    ```

- **Delete a product:**
    ```bash
    curl -X DELETE http://localhost:3000/products?id=5,6,7
    ```
    ```json
    {
      "result": [
          5,
          6,
          7
        ],
      "status": 200
    }
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
