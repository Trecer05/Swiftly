CREATE TABLE IF NOT EXISTS "users" (
	"id" serial NOT NULL,
	"email" varchar(100) NOT NULL UNIQUE,
	"number" integer NOT NULL UNIQUE,
	"password_hash" text NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "user_tokens" (
	"id" serial PRIMARY KEY,
	"user_id" integer NOT NULL,
	"refresh" text NOT NULL,
	"expired_at" timestamp with time zone NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	PRIMARY KEY ("user_id")
);

CREATE TABLE IF NOT EXISTS "expired_tokens" (
	"refresh" text NOT NULL,
	"user_id" integer NOT NULL
);

ALTER TABLE "users" ADD CONSTRAINT "users_fk0" FOREIGN KEY ("id") REFERENCES "user_tokens"("user_id");
ALTER TABLE "expired_tokens" ADD CONSTRAINT "expired_tokens_fk0" FOREIGN KEY ("token_id") REFERENCES "user_tokens"("id");