DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type WHERE typname = 'user_role'
	) THEN
		CREATE TYPE user_role 
		AS 
		ENUM('customer', 'admin');
	END IF;
END;
$$;
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type WHERE typname = 'user_status'
	) THEN
		CREATE TYPE "user_status" 
		AS 
		ENUM('active', 'blocked');
	END IF;
END;
$$;
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