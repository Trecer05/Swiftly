CREATE TABLE IF NOT EXISTS "users" (
	"id" integer NOT NULL,
	"email" varchar(100) NOT NULL UNIQUE,
	"number" integer NOT NULL UNIQUE,
	"password_hash" text NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "user_tokens" (
	"user_id" integer NOT NULL,
	"refresh" text NOT NULL,
	"expired_at" timestamp with time zone NOT NULL,
	PRIMARY KEY ("user_id")
);

ALTER TABLE "users" ADD CONSTRAINT "users_fk0" FOREIGN KEY ("id") REFERENCES "user_tokens"("user_id");