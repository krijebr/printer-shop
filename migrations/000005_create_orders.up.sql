CREATE TYPE order_status 
AS 
ENUM('new', 'in_progress', 'done');
CREATE TABLE IF NOT EXISTS "orders" (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	status order_status NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id)
);
ALTER TABLE orders ADD CONSTRAINT fk_orders FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE;
CREATE TABLE IF NOT EXISTS "order_products" (
	product_id uuid NOT NULL,
	order_id uuid NOT NULL,
	product_count integer NOT NULL,
	product_price float4 NOT NULL
);
ALTER TABLE order_products ADD CONSTRAINT fk_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE order_products ADD CONSTRAINT fk_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE;