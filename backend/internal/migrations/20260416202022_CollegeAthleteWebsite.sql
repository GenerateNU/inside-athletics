-- Modify "colleges" table
ALTER TABLE "public"."colleges" ADD COLUMN "athletics_website" character varying(500) NOT NULL;
-- Create "athletes" table
CREATE TABLE "public"."athletes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(200) NOT NULL,
  "sport_id" uuid NULL,
  "college_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_athletes_college" FOREIGN KEY ("college_id") REFERENCES "public"."colleges" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_athletes_sport" FOREIGN KEY ("sport_id") REFERENCES "public"."sports" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
