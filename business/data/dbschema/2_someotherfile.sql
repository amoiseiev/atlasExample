-- Add new schema named "myschema2"
CREATE SCHEMA "myschema2";
-- create "accounts2" table
CREATE TABLE "myschema2"."accounts2" (
                          "user_id" serial PRIMARY KEY,
                          "username" VARCHAR ( 50 ) UNIQUE NOT NULL,
                          "password" VARCHAR ( 50 ) NOT NULL,
                          "email" VARCHAR ( 255 ) UNIQUE NOT NULL,
                          "created_on" TIMESTAMP NOT NULL,
                          "last_login" TIMESTAMP NOT NULL
);