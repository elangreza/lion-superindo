BEGIN
;

CREATE SEQUENCE productseq START 100;

CREATE TABLE IF NOT EXISTS "products" (
    "id" BIGINT DEFAULT nextval('productseq') PRIMARY KEY,
    "name" VARCHAR UNIQUE NOT NULL,
    "quantity" BIGINT CHECK(quantity >= 0),
    "price" BIGINT CHECK(price >= 0),
    "type_name" VARCHAR REFERENCES types("name"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NULL
);

CREATE TRIGGER "log_product_update" BEFORE
UPDATE
    ON "products" FOR EACH ROW EXECUTE PROCEDURE log_update_master();

COMMIT;