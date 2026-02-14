-- Create "user_roles" join table
CREATE TABLE "public"."user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "role_id"),
  CONSTRAINT "user_roles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Backfill existing users with their current role
INSERT INTO "public"."user_roles" ("user_id", "role_id")
SELECT u."id", u."role_id"
FROM "public"."users" u
WHERE u."role_id" IS NOT NULL
ON CONFLICT DO NOTHING;

-- Remove legacy single-role column
ALTER TABLE "public"."users" DROP CONSTRAINT IF EXISTS "users_role_id_fkey";
DROP INDEX IF EXISTS public.users_role_id_idx;
ALTER TABLE "public"."users" DROP COLUMN IF EXISTS "role_id";
