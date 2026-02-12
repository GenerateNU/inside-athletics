-- Create "colleges" table
CREATE TABLE "public"."colleges" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(200) NOT NULL,
  "state" character varying(100) NOT NULL,
  "city" character varying(100) NOT NULL,
  "website" character varying(500) NULL,
  "academic_rank" smallint NULL,
  "division_rank" bigint NOT NULL,
  "logo" character varying(500) NULL,
  PRIMARY KEY ("id")
);
