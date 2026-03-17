-- Create "sport_follows" table
CREATE TABLE "public"."sport_follows" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "sport_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sport_follows_sport" FOREIGN KEY ("sport_id") REFERENCES "public"."sports" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_sport_follows_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
