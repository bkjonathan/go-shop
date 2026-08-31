-- reverse: create index "idx_refresh_tokens_token" to table: "refresh_tokens"
DROP INDEX "idx_refresh_tokens_token";
-- reverse: create index "idx_refresh_tokens_deleted_at" to table: "refresh_tokens"
DROP INDEX "idx_refresh_tokens_deleted_at";
-- reverse: create "refresh_tokens" table
DROP TABLE "refresh_tokens";
-- reverse: create index "idx_product_images_deleted_at" to table: "product_images"
DROP INDEX "idx_product_images_deleted_at";
-- reverse: create "product_images" table
DROP TABLE "product_images";
-- reverse: create index "idx_order_items_deleted_at" to table: "order_items"
DROP INDEX "idx_order_items_deleted_at";
-- reverse: create "order_items" table
DROP TABLE "order_items";
-- reverse: create index "idx_orders_deleted_at" to table: "orders"
DROP INDEX "idx_orders_deleted_at";
-- reverse: create "orders" table
DROP TABLE "orders";
-- reverse: create index "idx_cart_items_deleted_at" to table: "cart_items"
DROP INDEX "idx_cart_items_deleted_at";
-- reverse: create "cart_items" table
DROP TABLE "cart_items";
-- reverse: create index "idx_products_deleted_at" to table: "products"
DROP INDEX "idx_products_deleted_at";
-- reverse: create "products" table
DROP TABLE "products";
-- reverse: create index "idx_categories_deleted_at" to table: "categories"
DROP INDEX "idx_categories_deleted_at";
-- reverse: create "categories" table
DROP TABLE "categories";
-- reverse: create index "idx_carts_deleted_at" to table: "carts"
DROP INDEX "idx_carts_deleted_at";
-- reverse: create "carts" table
DROP TABLE "carts";
-- reverse: create index "idx_users_email" to table: "users"
DROP INDEX "idx_users_email";
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX "idx_users_deleted_at";
-- reverse: create "users" table
DROP TABLE "users";
