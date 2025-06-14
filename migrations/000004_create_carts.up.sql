
CREATE TABLE IF NOT EXISTS "carts" (
	user_id uuid NOT NULL,
	product_id uuid NOT NULL,
	count integer NOT NULL
);
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_carts_users'
	) THEN
			EXECUTE 'ALTER TABLE carts ADD CONSTRAINT fk_carts_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE';
	END IF;
END;
$$;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_carts_products'
	) THEN
		EXECUTE 'ALTER TABLE carts ADD CONSTRAINT fk_carts_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE';
	END IF;
END;
$$;