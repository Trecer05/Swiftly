CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	email varchar(100) NOT NULL UNIQUE,
	number varchar(100) NOT NULL UNIQUE,
	password_hash text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
	id serial PRIMARY KEY,
	user_id integer NOT NULL UNIQUE,
	token text NOT NULL,
	expired_at timestamp with time zone NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expired_tokens (
	user_id integer NOT NULL,
	token text NOT NULL
);

ALTER TABLE "user_tokens" ADD CONSTRAINT "user_tokens_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("id");
ALTER TABLE "expired_tokens" ADD CONSTRAINT "expired_tokens_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("id");