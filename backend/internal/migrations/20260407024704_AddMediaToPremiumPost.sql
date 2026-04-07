-- Modify "premium_posts" table
ALTER TABLE "public"."premium_posts" DROP COLUMN "attachment_key", DROP COLUMN "attachment_type", ADD COLUMN "media_id" uuid NULL, ADD CONSTRAINT "fk_premium_posts_media" FOREIGN KEY ("media_id") REFERENCES "public"."media" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Drop index "idx_tag_posts_postable" from table: "tag_posts"
DROP INDEX "public"."idx_tag_posts_postable";
-- Modify "tag_posts" table
ALTER TABLE "public"."tag_posts" ADD CONSTRAINT "fk_tag_posts_post" FOREIGN KEY ("postable_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_tag_posts_premium_post" FOREIGN KEY ("postable_id") REFERENCES "public"."premium_posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
