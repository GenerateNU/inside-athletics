-- Create "media" table
CREATE TABLE "public"."media" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "s3_key" character varying(200) NOT NULL,
  "title" character varying(200) NOT NULL,
  "media_type" character varying(200) NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
