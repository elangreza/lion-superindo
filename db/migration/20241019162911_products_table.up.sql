BEGIN
;

CREATE SEQUENCE productseq START 100;

CREATE TABLE IF NOT EXISTS "products" (
    "id" BIGINT DEFAULT nextval('productseq') PRIMARY KEY,
    "name" VARCHAR UNIQUE NOT NULL,
    "price" BIGINT CHECK(price >= 0),
    "product_type_name" VARCHAR REFERENCES product_types("name"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NULL
);

CREATE TRIGGER "log_products_update" BEFORE
UPDATE
    ON "products" FOR EACH ROW EXECUTE PROCEDURE log_update_master();

COMMIT;