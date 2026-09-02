-- modify "products" table
ALTER TABLE "products" ADD COLUMN "sku" text NOT NULL;
-- create index "idx_products_sku" to table: "products"
CREATE UNIQUE INDEX "idx_products_sku" ON "products" ("sku");
