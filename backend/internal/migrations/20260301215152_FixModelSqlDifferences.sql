-- Modify "permissions" table
ALTER TABLE "public"."permissions" DROP CONSTRAINT "permissions_action_resource_key";
-- Create index "permissions_action_resource_key" to table: "permissions"
CREATE UNIQUE INDEX "permissions_action_resource_key" ON "public"."permissions" ("action", "resource");
-- Create index "idx_permissions_deleted_at" to table: "permissions"
CREATE INDEX "idx_permissions_deleted_at" ON "public"."permissions" ("deleted_at");
-- Modify "roles" table
ALTER TABLE "public"."roles" DROP CONSTRAINT "roles_name_key";
-- Create index "roles_name_key" to table: "roles"
CREATE UNIQUE INDEX "roles_name_key" ON "public"."roles" ("name");
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deleted_at");
-- Modify "role_permissions" table
ALTER TABLE "public"."role_permissions" DROP CONSTRAINT "role_permissions_permission_id_fkey", DROP CONSTRAINT "role_permissions_role_id_fkey", ADD CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_role_permissions_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "user_roles" table
ALTER TABLE "public"."user_roles" DROP CONSTRAINT "user_roles_role_id_fkey", DROP CONSTRAINT "user_roles_user_id_fkey", ADD CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
