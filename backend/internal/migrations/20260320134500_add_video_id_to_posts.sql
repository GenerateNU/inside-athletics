-- Add "video_id" column to "posts"
ALTER TABLE "public"."posts"
  ADD COLUMN "video_id" uuid NULL;

-- Add foreign key constraint
ALTER TABLE "public"."posts"
  ADD CONSTRAINT "posts_video_id_fkey"
  FOREIGN KEY ("video_id") REFERENCES "public"."videos" ("id")
  ON UPDATE NO ACTION
  ON DELETE SET NULL;

-- Optional index (recommended)
CREATE INDEX "posts_video_id_idx" ON "public"."posts" ("video_id");