-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "expected_grad_year" TYPE bigint, ALTER COLUMN "division" TYPE bigint;
-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "author_id" text NULL,
  "sport_id" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" character varying(100) NOT NULL,
  "content" character varying(5000) NOT NULL,
  "likes" bigint NULL,
  "is_anonymous" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_posts_deleted_at" to table: "posts"
CREATE INDEX "idx_posts_deleted_at" ON "public"."posts" ("deleted_at");
