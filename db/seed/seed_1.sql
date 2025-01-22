INSERT INTO
    public."product_types" ("name")
VALUES
    ('sayuran'),
    ('protein'),
    ('buah'),
    ('snack');

INSERT INTO
    public .products (
        "name",
        quantity,
        price,
        product_type_name
    )
VALUES
    ('Sawi', 10, 3000, 'sayuran'),
    ('Kangkung', 5, 2000, 'sayuran'),
    ('Tauge', 1, 1000, 'sayuran'),
    ('Tempe', 2, 9000, 'protein'),
    ('Pepaya', 4, 4000, 'buah'),
    ('Singkong', 3, 5000, 'buah'),
    ('Donat', 7, 6000, 'snack');