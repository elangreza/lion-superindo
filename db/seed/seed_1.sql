INSERT INTO
    public."types" ("name")
VALUES
    ('Sayuran'),
    ('Protein'),
    ('Buah'),
    ('Snack');

INSERT INTO
    public .products (
        "name",
        quantity,
        price,
        type_name
    )
VALUES
    ('Sawi', 10, 3000, 'Sayuran'),
    ('Kangkung', 5, 2000, 'Sayuran'),
    ('Tauge', 1, 1000, 'Sayuran'),
    ('Tempe', 2, 9000, 'Protein'),
    ('Pepaya', 4, 4000, 'Buah'),
    ('Singkong', 3, 5000, 'Buah'),
    ('Donat', 7, 6000, 'Snack');