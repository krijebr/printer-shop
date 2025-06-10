DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type WHERE typname = 'order_status'
	) THEN
		CREATE TYPE order_status 
		AS 
		ENUM('new', 'in_progress', 'done');
	END IF;
END;
$$;
CREATE TABLE IF NOT EXISTS "orders" (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	status order_status NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id)
);
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_orders_users'
	) THEN
		EXECUTE 'ALTER TABLE orders ADD CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE';
	END IF;
END;
$$;

CREATE TABLE IF NOT EXISTS "order_products" (
	product_id uuid NOT NULL,
	order_id uuid NOT NULL,
	product_count integer NOT NULL,
	product_price float4 NOT NULL
);
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_order_products_products'
	) THEN
		EXECUTE 'ALTER TABLE order_products ADD CONSTRAINT fk_order_products_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE';
	END IF;
END;
$$;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_order_products_orders'
	) THEN		
		EXECUTE 'ALTER TABLE order_products ADD CONSTRAINT fk_order_products_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE';
	END IF;
END;
$$;