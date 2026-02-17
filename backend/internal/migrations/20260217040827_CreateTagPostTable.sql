-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "expected_grad_year" TYPE bigint, ALTER COLUMN "division" TYPE bigint;
-- Create "tag_posts" table
CREATE TABLE "public"."tag_posts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "post_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "tag_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id", "post_id", "tag_id")
);
