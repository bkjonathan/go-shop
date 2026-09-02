-- reverse: create index "idx_products_sku" to table: "products"
DROP INDEX "idx_products_sku";
-- reverse: modify "products" table
ALTER TABLE "products" DROP COLUMN "sku";
