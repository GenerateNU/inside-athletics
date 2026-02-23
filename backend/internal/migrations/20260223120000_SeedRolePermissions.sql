-- Seed permissions for posts, comments, and likes
INSERT INTO "public"."permissions" ("action", "resource") VALUES
  ('get', 'role'),
  ('get', 'permission'),
  ('create', 'post'),
  ('update', 'post'),
  ('delete', 'post'),
  ('update_own', 'post'),
  ('delete_own', 'post'),
  ('create', 'comment'),
  ('update', 'comment'),
  ('delete', 'comment'),
  ('update_own', 'comment'),
  ('delete_own', 'comment'),
  ('create', 'like'),
  ('delete', 'like')
ON CONFLICT DO NOTHING;

-- Assign permissions for regular users
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON (
  (p."resource" = 'post' AND p."action" IN ('create', 'update_own', 'delete_own')) OR
  (p."resource" = 'comment' AND p."action" IN ('create', 'update_own', 'delete_own')) OR
  (p."resource" = 'like' AND p."action" IN ('create', 'delete'))
)
WHERE r."name" = 'user'
ON CONFLICT DO NOTHING;

-- Assign permissions for moderators (user perms + edit/delete others' posts/comments)
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON (
  (p."resource" = 'post' AND p."action" IN ('create', 'update_own', 'delete_own', 'update', 'delete')) OR
  (p."resource" = 'comment' AND p."action" IN ('create', 'update_own', 'delete_own', 'update', 'delete')) OR
  (p."resource" = 'like' AND p."action" IN ('create', 'delete'))
)
WHERE r."name" = 'moderator'
ON CONFLICT DO NOTHING;

-- Admins should have access to everything
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p ON TRUE
WHERE r."name" = 'admin'
ON CONFLICT DO NOTHING;
