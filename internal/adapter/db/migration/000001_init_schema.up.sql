CREATE TABLE "product" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) UNIQUE
);

CREATE TABLE "category" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) UNIQUE
);

CREATE TABLE "product_category" (
    "product_id" bigserial,
    "category_id" bigserial,
    UNIQUE ("product_id", "category_id")
);

CREATE TABLE "user" (
    "id" bigserial PRIMARY KEY,
    "login" varchar(255) UNIQUE,
    "password" varchar
);

CREATE TABLE "session" (
    "id" bigserial PRIMARY KEY,
    "token" varchar UNIQUE,
    "user_id" bigserial,
    "expiry" timestamp
);

ALTER TABLE "product_category" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE CASCADE;

ALTER TABLE "product_category" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id") ON DELETE CASCADE;

ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");