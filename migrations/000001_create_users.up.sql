CREATE TYPE user_role 
AS 
ENUM('customer', 'admin');
CREATE TYPE "user_status" 
AS 
ENUM('active', 'blocked');
CREATE TABLE IF NOT EXISTS "users" (
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