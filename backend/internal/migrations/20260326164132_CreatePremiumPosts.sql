-- Create "premium_posts" table
CREATE TABLE "public"."premium_posts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "author_id" uuid NOT NULL,
  "sport_id" uuid NULL,
  "college_id" uuid NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" character varying(100) NOT NULL,
  "content" character varying(5000) NOT NULL,
  "attachment_key" text NULL,
  "attachment_type" character varying(10) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_premium_posts_author" FOREIGN KEY ("author_id") REFERENCES "public"."users"("id") ON DELETE CASCADE,
  CONSTRAINT "fk_premium_posts_sport" FOREIGN KEY ("sport_id") REFERENCES "public"."sports"("id") ON DELETE SET NULL,
  CONSTRAINT "fk_premium_posts_college" FOREIGN KEY ("college_id") REFERENCES "public"."colleges"("id") ON DELETE SET NULL
);
-- Create index "idx_premium_posts_deleted_at" to table: "premium_posts"
CREATE INDEX "idx_premium_posts_deleted_at" ON "public"."premium_posts" ("deleted_at");

-- Create "tag_posts" join table for premium_posts <-> tags
CREATE TABLE "public"."tag_posts" (
  "tag_id" uuid NOT NULL,
  "premium_post_id" uuid NOT NULL,
  PRIMARY KEY ("tag_id", "premium_post_id"),
  CONSTRAINT "fk_tag_posts_tag" FOREIGN KEY ("tag_id") REFERENCES "public"."tags"("id") ON DELETE CASCADE,
  CONSTRAINT "fk_tag_posts_premium_post" FOREIGN KEY ("premium_post_id") REFERENCES "public"."premium_posts"("id") ON DELETE CASCADE
);