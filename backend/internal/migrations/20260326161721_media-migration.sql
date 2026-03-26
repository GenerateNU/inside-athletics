-- Modify "posts" table
ALTER TABLE "public"."posts" DROP COLUMN "video_id";
-- Create "media" table
CREATE TABLE "public"."media" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "s3_key" text NULL,
  "title" text NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Drop "videos" table
DROP TABLE "public"."videos";
