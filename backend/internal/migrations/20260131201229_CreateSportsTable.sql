-- Create "sports" table
CREATE TABLE "public"."sports" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(100) NOT NULL,
  "popularity" integer NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sports_deleted_at" to table: "sports"
CREATE INDEX "idx_sports_deleted_at" ON "public"."sports" ("deleted_at");
