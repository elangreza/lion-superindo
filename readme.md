# Coding test Lion Superindo:

Buatlah sebuah API dengan endpoint /product untuk menambahkan dan mengambil data product super
indo, dengan spesifikasi sebagai berikut :

- [x] Dapat melakukan penambahkan data product
- [x] Dapat menampilkan list data product
- [x] Dapat melakukan pencarian bedasarkan nama dan id product
- [x] Dapat melakukan filter produk berdasarkan tipe produk Sayuran, Protein, Buah dan Snack
- [x] Dapat melakukan sorting berdasarkan tanggal, harga dan nama product

## Tech Stack :

- [x] Language : Golang
- [x] Database : SQL / NoSQL + Seeder + migration
- [x] Cache : Redis
- [x] Dependency Injection : wire (Optional)
- [x] Unittest (Optional)
- [x] Docker (Optional)

## How to Run the Application

1. **Start the application:**

   ```sh
   make up
   ```

   This will build the application with Docker Compose and automatically seed the database.

2. **Add a Product:**

   ```sh
   curl --location 'http://localhost:8080/product' \
     --header 'Content-Type: application/json' \
     --data '{
       "name": "kopi luwak",
       "type": "Snack",
       "price": 10000
     }'
   ```

   See the [POST `/product` API documentation](#post-product) or [Swagger API](http://localhost:8080/swagger/index.html#/product/post_product) for details.

3. **Retrieve Products:**

   ```sh
   curl --location 'http://localhost:8080/product?page=1&search=pe&type=buah&type=proteinas&sort=id:asc'
   ```

   See the [GET `/product` API documentation](#get-product) or [Swagger API](http://localhost:8080/swagger/index.html#/product/get_product) for details.

4. **API Documentation:**

   See the [API documentation section](#api-documentation) for docs

   or

   see the [Swagger UI](http://localhost:8080/swagger/index.html) for trying api directly.

5. **Shutdown the stack:**

   ```sh
   make down
   ```

## API Documentation

### `/product` Endpoint

Only **POST** and **GET** methods are supported for this endpoint. Any other HTTP method will return a "Method Not Allowed" (405) error.

#### POST `/product`

- **Purpose:** Create a new product.
- **Request Body:**
  ```json
  {
    "name": "kopi luwak",
    "type": "Snack",
    "price": 10000
  }
  ```
- **Responses:**
  - **201 Created**
    ```json
    { "data": { "id": 168 } }
    ```
  - **409 Conflict** (Product already exists)
    ```json
    { "error": "product already exist" }
    ```

#### GET `/product`

- **Purpose:** Retrieve a list of products.
- **Query Parameters:**

  - `search` — Search by id or name.  
    _Example:_ `/product?search=semangka`
  - `sort` — Sort by `id`, `name`, `price`, or `created_at`.  
    _Format:_ `key:asc` or `key:desc`  
    _Example:_ `/product?sort=created_at:asc&sort=name:desc&sort=price:asc`
  - `type` — Filter by product type.  
    _Example:_ `/product?type=buah&type=snack`
  - `page` — Pagination.  
    _Example:_ `/product?page=1`
  - `limit` — Items per page.  
    _Example:_ `/product?limit=10`

- **Response:**
  - **200 OK**
    ```json
    {
      "data": {
        "total_data": 2,
        "total_page": 2,
        "products": [
          {
            "id": 168,
            "name": "kopi luwak",
            "price": 10000,
            "type": "snack",
            "created_at": "2025-01-23T10:51:05.445274Z"
          },
          {
            "id": 167,
            "name": "kopi Arabica",
            "price": 10000,
            "type": "snack",
            "created_at": "2025-01-23T10:39:33.187086Z"
          }
        ]
      }
    }
    ```

#### Other Methods

- **405 Method Not Allowed**
  ```json
  { "error": "invalid method" }
  ```
