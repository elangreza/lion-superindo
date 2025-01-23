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

## Documentation

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

            | query name | type             | default value | example                                                       |
            | ---------- | ---------------- | ------------- | ------------------------------------------------------------- |
            | search     | string           | empty         | /product?search=semangka                                      |
            | sorts      | array of string  | empty         | /product?sorts=updated_at:asc&sorts=name:desc&sorts=price:asc |
            | types      | array of string  | empty         | /product?types=buah&types=snack                               |
            | page       | positive integer | 1             | /product?page=1                                               |
            | limit      | positive integer | 5             | /product?limit=10                                             |

        - another method will be rejected with status **405**

          ```json
          {
            "error": "invalid method"
          }
          ```
