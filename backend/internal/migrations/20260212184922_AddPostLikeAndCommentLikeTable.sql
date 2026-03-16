-- Create "comment_likes" table
CREATE TABLE "public"."comment_likes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "comment_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_comment_likes_comment" FOREIGN KEY ("comment_id") REFERENCES "public"."comments" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_comment_likes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_comment_likes_deleted_at" to table: "comment_likes"
CREATE INDEX "idx_comment_likes_deleted_at" ON "public"."comment_likes" ("deleted_at");
-- Create "post_likes" table
CREATE TABLE "public"."post_likes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_likes_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_post_likes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_post_likes_deleted_at" to table: "post_likes"
CREATE INDEX "idx_post_likes_deleted_at" ON "public"."post_likes" ("deleted_at");
