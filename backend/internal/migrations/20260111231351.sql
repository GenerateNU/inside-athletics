-- Create "users" table
CREATE TABLE `users` (
  `id` bigserial NULL,
  `created_at` timestamptz NULL,
  `updated_at` timestamptz NULL,
  `deleted_at` timestamptz NULL,
  `name` varchar NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
