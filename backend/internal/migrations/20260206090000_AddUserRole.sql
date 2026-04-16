-- Add role_id to users
ALTER TABLE "public"."users"
  ADD COLUMN "role_id" uuid NULL;

-- Backfill existing users to the "user" role
UPDATE "public"."users" u
SET "role_id" = r."id"
FROM "public"."roles" r
WHERE r."name" = 'user'
  AND u."role_id" IS NULL;

-- Enforce not-null and add FK + index
ALTER TABLE "public"."users"
  ALTER COLUMN "role_id" SET NOT NULL;

ALTER TABLE "public"."users"
  ADD CONSTRAINT "users_role_id_fkey"
  FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id")
  ON UPDATE NO ACTION ON DELETE RESTRICT;

CREATE INDEX "users_role_id_idx" ON "public"."users" ("role_id");
