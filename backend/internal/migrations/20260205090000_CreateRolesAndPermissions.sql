-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(50) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "roles_name_key" UNIQUE ("name")
);

-- Create "permissions" table
CREATE TABLE "public"."permissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "action" character varying(50) NOT NULL,
  "resource" character varying(50) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "permissions_action_resource_key" UNIQUE ("action", "resource")
);

-- Create "role_permissions" join table
CREATE TABLE "public"."role_permissions" (
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("role_id", "permission_id"),
  CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Seed base roles
INSERT INTO "public"."roles" ("name") VALUES
  ('user'),
  ('admin'),
  ('moderator');

-- Seed admin permissions (roles + permissions management)
INSERT INTO "public"."permissions" ("action", "resource") VALUES
  ('create', 'role'),
  ('update', 'role'),
  ('delete', 'role'),
  ('create', 'permission'),
  ('update', 'permission'),
  ('delete', 'permission');

-- Assign admin role permissions
INSERT INTO "public"."role_permissions" ("role_id", "permission_id")
SELECT r."id", p."id"
FROM "public"."roles" r
JOIN "public"."permissions" p
  ON p."resource" IN ('role', 'permission')
WHERE r."name" = 'admin';
