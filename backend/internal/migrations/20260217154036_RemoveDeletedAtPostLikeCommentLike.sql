-- Modify "comment_likes" table
ALTER TABLE "public"."comment_likes" DROP COLUMN "deleted_at";
-- Modify "post_likes" table
ALTER TABLE "public"."post_likes" DROP COLUMN "deleted_at";
