-- Create "college_follows" table
CREATE TABLE "public"."college_follows" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "college_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_college_follows_college" FOREIGN KEY ("college_id") REFERENCES "public"."colleges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_college_follows_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
