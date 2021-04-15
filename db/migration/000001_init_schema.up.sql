CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "title" varchar NOT NULL,
  "content" varchar NOT NULL,
  "votes" int NOT NULL DEFAULT 0,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "votes" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "post_id" int NOT NULL,
  "voted" bool NOT NULL
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "votes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "votes" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");

CREATE INDEX ON "posts" ("user_id");

CREATE INDEX ON "votes" ("user_id");

CREATE INDEX ON "votes" ("post_id");

CREATE UNIQUE INDEX "user_post_index" ON "votes" ("user_id", "post_id");

COMMENT ON COLUMN "votes"."voted" IS 'true if liked and false if disliked. No rows if no response';
