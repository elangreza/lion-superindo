# Coding test SuperIndo:

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
- [ ] Dependency Injection : wire (Optional)
- [x] Unittest (Optional)
- [x] Docker (Optional)

## How to Run the Application

1. Run the application:

   ```sh
   make up
   ```

2. Try adding a product:

   ```curl
   curl --location 'http://localhost:8080/product' \
   --header 'Content-Type: application/json' \
   --data '{
       "name":"kopi luwak",
       "type":"Snack",
       "price":10000
   }'
   ```

3. Try retrieving products:

   ```curl
  curl --location 'localhost:8080/product?page=1&search=pe&type=buah&type=proteinas&sort=id%3Aasc'
   ```

## API Documentation

1.  `/product` API

    - The **POST** method is used to create a new product. The payload must follow this format:

      ```json
      {
        "name": "kopi luwak",
        "type": "Snack",
        "price": 10000
      }
      ```

      If the product is created, you will receive status **201**:

      ```json
      {
        "data": {
          "id": 168
        }
      }
      ```

      If the product already exists in the database, you will receive status **409**:

      ```json
      {
        "error": "product already exist"
      }
      ```

    - The **GET** method is used to retrieve a list of products. This method can be combined with query parameters, with 5 possible values:

      1. **search**
         Used to search by id or name.

         Example:

         ```
         /product?search=semangka
         ```

      2. **sort**
         Used to sort the data; it can only be used for _name_, _price_, and _created_at_. The format for sorting is `key:asc` or `key:desc`. Example:

         ```
         /product?sort=created_at:asc&sort=name:desc&sort=price:asc
         ```

      3. **type**
         Used to filter the data by product type. Example:

         ```
         /product?type=buah&type=snack
         ```

      4. **page**
         Used for pagination. Example:

         ```
         /product?page=1
         ```

      5. **limit**
         Used to limit the number of items per page. Example:

         ```
         /product?page=1
         ```

      The response will look like this, with status **200**:

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

    - Any other method will be rejected with status **405**:

      ```json
      {
        "error": "invalid method"
      }
      ```
