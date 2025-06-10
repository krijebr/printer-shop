CREATE TABLE IF NOT EXISTS "producers" (
	id uuid NOT NULL,
	name varchar NOT NULL,
	description varchar NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT producers_pk PRIMARY KEY (id)
);