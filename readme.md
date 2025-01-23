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

## How to run the Application

1. copy env file

   ```sh
   cp ./env.example .env
   ```

2. build docker file

   ```sh
   make up
   ```

3. run the seed migration

   ```sh
   cat ./db/seed/seed_1.sql | docker exec -i superindo-database psql -h localhost -U superindo -f-
   ```

4. try with adding product

   ```curl
   curl --location 'http://localhost:8080/product' \
   --header 'Content-Type: application/json' \
   --data '{
       "name":"kopi luwak",
       "type":"Snack",
       "price":10000
   }'
   ```

5. try with getting products

   ```curl
   curl --location 'http://localhost:8080/product?sorts=updated_at%3Aasc&sorts=name%3Adesc&page=1&limit=10&sorts=price%3Aasc'
   ```

## API Documentation

1.  `/product` API

    - method **POST** is used for creating new product. the payload must followed this format

      ```json
      {
        "name": "kopi luwak",
        "type": "Snack",
        "price": 10000
      }
      ```

      if the product is created will get status **201**

      ```json
      {
        "data": {
          "id": 168
        }
      }
      ```

      if the product is exist in db with status **409**

      ```json
      {
        "error": "product already exist"
      }
      ```

    - method **GET** is used for retrieving list of product. This method can combined with query params with 5 possible values

      1. **search**
         is used for searching by id or name

         example

         ```
         /product?search=semangka
         ```

      2. **sorts**
         is used for sorting the data, can be used only for _name_, _price_, and _updated_at_. The format of sorts is `key` + ":" + `asc` or `key` + ":" + `desc`. Take a look at this example

         ```
         /product?sorts=updated_at:asc&sorts=name:desc&sorts=price:asc
         ```

      3. **types**
         is used for filtering the data, can be used with the value of product types. Take a look at this example

         ```
         /product?types=buah&types=snack
         ```

      4. **page**
         is used for splitting the data with page. Take a look at this example

         ```
         /product?page=1
         ```

      5. **limit**
         is used for limiting the data each page. Take a look at this example

         ```
         /product?page=1
         ```

      and the response will be like this and status is **200**

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
              "updated_at": "2025-01-23T10:51:05.445274Z"
            },
            {
              "id": 167,
              "name": "kopi Arabica",
              "price": 10000,
              "type": "snack",
              "updated_at": "2025-01-23T10:39:33.187086Z"
            }
          ]
        }
      }
      ```

    - another method will be rejected with status **405**

      ```json
      {
        "error": "invalid method"
      }
      ```
