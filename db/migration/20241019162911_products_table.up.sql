BEGIN
;

CREATE SEQUENCE productseq START 100;

CREATE TABLE IF NOT EXISTS "products" (
    "id" BIGINT DEFAULT nextval('productseq') PRIMARY KEY,
    "name" VARCHAR UNIQUE NOT NULL,
    "price" BIGINT CHECK(price >= 0),
    "product_type_name" VARCHAR REFERENCES product_types("name"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;