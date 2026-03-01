-- Drop index "idx_comment_likes_user_id_comment_id" from table: "comment_likes"
DROP INDEX "public"."idx_comment_likes_user_id_comment_id";
-- Modify "permissions" table
ALTER TABLE "public"."permissions" DROP CONSTRAINT "permissions_action_resource_key";
-- Create index "idx_permissions_deleted_at" to table: "permissions"
CREATE INDEX "idx_permissions_deleted_at" ON "public"."permissions" ("deleted_at");
-- Drop index "idx_post_likes_user_id_post_id" from table: "post_likes"
DROP INDEX "public"."idx_post_likes_user_id_post_id";
-- Modify "role_permissions" table
ALTER TABLE "public"."role_permissions" DROP CONSTRAINT "role_permissions_permission_id_fkey", DROP CONSTRAINT "role_permissions_role_id_fkey";
-- Modify "roles" table
ALTER TABLE "public"."roles" DROP CONSTRAINT "roles_name_key";
-- Create index "roles_name_key" to table: "roles"
CREATE UNIQUE INDEX "roles_name_key" ON "public"."roles" ("name");
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deleted_at");
-- Modify "user_roles" table
ALTER TABLE "public"."user_roles" DROP CONSTRAINT "user_roles_role_id_fkey", DROP CONSTRAINT "user_roles_user_id_fkey";
