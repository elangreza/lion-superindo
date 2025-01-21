BEGIN
;

CREATE TABLE IF NOT EXISTS "product_types" (
    "name" VARCHAR PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NULL
);

CREATE TRIGGER "log_product_types_update" BEFORE
UPDATE
    ON "product_types" FOR EACH ROW EXECUTE PROCEDURE log_update_master();

COMMIT;