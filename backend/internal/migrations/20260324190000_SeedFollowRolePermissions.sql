-- Seed permissions for follow resources
INSERT INTO "public"."permissions" ("action", "resource") VALUES
  ('create', 'tagfollow'),
  ('delete', 'tagfollow'),
  ('delete_own', 'tagfollow'),
  ('create', 'sportfollow'),
  ('delete', 'sportfollow'),
  ('delete_own', 'sportfollow'),
  ('create', 'collegefollow'),
  ('delete', 'collegefollow'),
  ('delete_own', 'collegefollow')
ON CONFLICT DO NOTHING;

-- Assign follow permissions for regular users
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON (
  (p."resource" = 'tagfollow' AND p."action" IN ('create', 'delete_own')) OR
  (p."resource" = 'sportfollow' AND p."action" IN ('create', 'delete_own')) OR
  (p."resource" = 'collegefollow' AND p."action" IN ('create', 'delete_own'))
)
WHERE r."name" = 'user'
ON CONFLICT DO NOTHING;

-- Assign follow permissions for moderators
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON (
  (p."resource" = 'tagfollow' AND p."action" IN ('create', 'delete_own')) OR
  (p."resource" = 'sportfollow' AND p."action" IN ('create', 'delete_own')) OR
  (p."resource" = 'collegefollow' AND p."action" IN ('create', 'delete_own'))
)
WHERE r."name" = 'moderator'
ON CONFLICT DO NOTHING;

-- Assign all follow permissions for admins
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON (
  p."resource" IN ('tagfollow', 'sportfollow', 'collegefollow')
)
WHERE r."name" = 'admin'
ON CONFLICT DO NOTHING;
