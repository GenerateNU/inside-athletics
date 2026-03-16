-- Add unique constraints for upsert (ON CONFLICT) on like tables
CREATE UNIQUE INDEX "idx_post_likes_user_id_post_id" ON "public"."post_likes" ("user_id", "post_id");
CREATE UNIQUE INDEX "idx_comment_likes_user_id_comment_id" ON "public"."comment_likes" ("user_id", "comment_id");
