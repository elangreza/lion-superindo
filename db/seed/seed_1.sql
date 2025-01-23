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
        price,
        product_type_name
    )
VALUES
    ('Sawi', 3000, 'sayuran'),
    ('Kangkung', 2000, 'sayuran'),
    ('Tauge', 1000, 'sayuran'),
    ('Tempe', 9000, 'protein'),
    ('Pepaya', 4000, 'buah'),
    ('Singkong', 5000, 'buah'),
    ('Donat', 6000, 'snack');