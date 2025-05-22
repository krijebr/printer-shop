CREATE TYPE user_role 
AS 
ENUM('customer', 'admin');
CREATE TYPE user_status 
AS 
ENUM('active', 'blocked');
CREATE TABLE users (
	id uuid NOT NULL,
	first_name varchar NOT NULL,
	last_name varchar NOT NULL,
	email varchar NOT NULL,
	password_hash varchar NOT NULL,
	status user_status NOT NULL,
	role user_role NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE producers (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	description varchar NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT producers_pk PRIMARY KEY (id)
);
CREATE TYPE product_status 
AS 
ENUM('published', 'hidden');
CREATE TABLE products (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	price float4 NOT NULL,
	producer_id uuid NOT NULL,
	status public.product_status NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT newtable_pk PRIMARY KEY (id)
);
ALTER TABLE public.products ADD CONSTRAINT fk_producers FOREIGN KEY (producer_id) REFERENCES public.producers(id) ON DELETE RESTRICT ON UPDATE CASCADE;
