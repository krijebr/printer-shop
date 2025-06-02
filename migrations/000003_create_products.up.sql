CREATE TYPE "product_status"
AS 
ENUM('published', 'hidden');
CREATE TABLE IF NOT EXISTS "products" (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	price float4 NOT NULL,
	producer_id uuid NOT NULL,
	status public.product_status NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT products_pk PRIMARY KEY (id)
);
ALTER TABLE public.products ADD CONSTRAINT fk_producers FOREIGN KEY (producer_id) REFERENCES public.producers(id) ON DELETE RESTRICT ON UPDATE CASCADE;