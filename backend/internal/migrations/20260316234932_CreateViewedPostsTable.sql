-- Create "viewed_posts" table
CREATE TABLE "public"."viewed_posts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_viewed_posts_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_viewed_posts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_viewed_posts_deleted_at" to table: "viewed_posts"
CREATE INDEX "idx_viewed_posts_deleted_at" ON "public"."viewed_posts" ("deleted_at");
-- Create index "idx_viewed_posts_user_post" to table: "viewed_posts"
CREATE UNIQUE INDEX "idx_viewed_posts_user_post" ON "public"."viewed_posts" ("user_id", "post_id");
