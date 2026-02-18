-- Create "users" table
CREATE TABLE `users` (
  `id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `email` text NOT NULL,
  `password_hash` text NOT NULL,
  `name` text NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
