DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type WHERE typname = 'product_status'
	) THEN
		CREATE TYPE "product_status"
		AS 
		ENUM('published', 'hidden');
	END IF;
END;
$$;
CREATE TABLE IF NOT EXISTS "products" (
	id uuid NOT NULL,
	name varchar NOT NULL,
	price float4 NOT NULL,
	producer_id uuid NOT NULL,
	status product_status NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT products_pk PRIMARY KEY (id)
);
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_products_producers'
	) THEN
		EXECUTE 'ALTER TABLE products ADD CONSTRAINT fk_products_producers FOREIGN KEY (producer_id) REFERENCES producers(id) ON DELETE RESTRICT ON UPDATE CASCADE';
	END IF;
END;
$$;