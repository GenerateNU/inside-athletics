-- Seed permissions for premium post resources
INSERT INTO "public"."permissions" ("action", "resource") VALUES
  ('create', 'premiumpost'),
  ('delete', 'premiumpost'),
  ('delete_own', 'premiumpost'),
  ('update_own', 'premiumpost'),
  ('update', 'premiumpost')
ON CONFLICT DO NOTHING;

-- Assign all premium post permissions to moderators and admins
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p
  ON p."resource" = 'premiumpost'
WHERE r."name" IN ('moderator', 'admin')
ON CONFLICT DO NOTHING;