-- Create premium_user role
INSERT INTO "public"."roles" ("name") VALUES ('premium_user')
ON CONFLICT DO NOTHING;

-- Assign all user permissions to premium_user
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."role_permissions" rp ON rp."role_id" = (
  SELECT "id" FROM "public"."roles" WHERE "name" = 'user'
)
JOIN "public"."permissions" p ON p."id" = rp."permission_id"
WHERE r."name" = 'premium_user'
ON CONFLICT DO NOTHING;

-- Grant premium_user access to get premiumpost
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p
  ON p."resource" = 'premiumpost' AND p."action" = 'get'
WHERE r."name" = 'premium_user'
ON CONFLICT DO NOTHING;
