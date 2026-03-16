-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "is_anonymous" boolean NOT NULL DEFAULT false,
  "parent_comment_id" uuid NULL,
  "post_id" uuid NOT NULL,
  "description" character varying(3000) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_comments_parent_comment" FOREIGN KEY ("parent_comment_id") REFERENCES "public"."comments" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_comments_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_comments_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
