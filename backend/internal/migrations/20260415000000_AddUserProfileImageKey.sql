-- Add profile image key to users
ALTER TABLE "public"."users"
ADD COLUMN "profile_image_key" character varying(200) NULL;
