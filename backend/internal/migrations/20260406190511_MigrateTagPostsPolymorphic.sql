-- Modify "premium_posts" table
ALTER TABLE "public"."premium_posts" ALTER COLUMN "author_id" DROP NOT NULL, ALTER COLUMN "attachment_type" SET DEFAULT NULL::character varying;
-- Modify "tag_posts" table
ALTER TABLE "public"."tag_posts" DROP CONSTRAINT "tag_posts_pkey", DROP CONSTRAINT "fk_tag_posts_post", DROP COLUMN "post_id", ALTER COLUMN "tag_id" DROP DEFAULT, ADD COLUMN "postable_id" uuid NOT NULL, ADD COLUMN "postable_type" character varying(20) NOT NULL, ADD PRIMARY KEY ("id");
-- Add index for polymorphic lookups
CREATE INDEX "idx_tag_posts_postable" ON "public"."tag_posts" ("postable_id", "postable_type");