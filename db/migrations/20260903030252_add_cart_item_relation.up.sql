-- modify "cart_items" table
ALTER TABLE "cart_items" DROP CONSTRAINT "fk_cart_items_cart", ADD CONSTRAINT "fk_carts_cart_items" FOREIGN KEY ("cart_id") REFERENCES "carts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
