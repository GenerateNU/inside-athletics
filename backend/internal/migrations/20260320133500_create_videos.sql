-- Create "videos" table
CREATE TABLE "public"."videos" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "s3_key" text NOT NULL,
  "title" text NOT NULL,
  PRIMARY KEY ("id")
);

-- Index for soft deletes (matches your sports table)
CREATE INDEX "idx_videos_deleted_at" ON "public"."videos" ("deleted_at");