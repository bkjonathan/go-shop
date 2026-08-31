-- create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "phone" text NULL,
  "is_active" boolean NULL DEFAULT true,
  "role" text NULL DEFAULT 'customer',
  PRIMARY KEY ("id")
);
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- create "carts" table
CREATE TABLE "carts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_cart" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_carts_deleted_at" to table: "carts"
CREATE INDEX "idx_carts_deleted_at" ON "carts" ("deleted_at");
-- create "categories" table
CREATE TABLE "categories" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "categories" ("deleted_at");
-- create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "category_id" bigint NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "price" numeric NOT NULL,
  "stock" bigint NULL DEFAULT 0,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");
-- create "cart_items" table
CREATE TABLE "cart_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "cart_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_cart_items_cart" FOREIGN KEY ("cart_id") REFERENCES "carts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_products_cart_items" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_cart_items_deleted_at" to table: "cart_items"
CREATE INDEX "idx_cart_items_deleted_at" ON "cart_items" ("deleted_at");
-- create "orders" table
CREATE TABLE "orders" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "total_amount" numeric NOT NULL,
  "status" text NULL DEFAULT 'pending',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_orders" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_orders_deleted_at" to table: "orders"
CREATE INDEX "idx_orders_deleted_at" ON "orders" ("deleted_at");
-- create "order_items" table
CREATE TABLE "order_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "price" numeric NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_order_items" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_products_order_items" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_order_items_deleted_at" to table: "order_items"
CREATE INDEX "idx_order_items_deleted_at" ON "order_items" ("deleted_at");
-- create "product_images" table
CREATE TABLE "product_images" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "product_id" bigint NOT NULL,
  "url" text NOT NULL,
  "alt_text" text NULL,
  "is_primary" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_images" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_product_images_deleted_at" to table: "product_images"
CREATE INDEX "idx_product_images_deleted_at" ON "product_images" ("deleted_at");
-- create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_refresh_token" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_refresh_tokens_deleted_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_deleted_at" ON "refresh_tokens" ("deleted_at");
-- create index "idx_refresh_tokens_token" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token" ON "refresh_tokens" ("token");
