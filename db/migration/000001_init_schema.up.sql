CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "title" varchar NOT NULL,
  "descr" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "likes" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "post_id" int NOT NULL,
  "liked" bool NOT NULL,
  "liked_at" timestamptz DEFAULT (now())
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "likes" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");

CREATE INDEX ON "posts" ("user_id");

CREATE INDEX ON "likes" ("user_id");

CREATE INDEX ON "likes" ("post_id");

CREATE INDEX ON "likes" ("user_id", "post_id");

COMMENT ON COLUMN "likes"."liked" IS 'true if liked and false if disliked. No rows if no response';
