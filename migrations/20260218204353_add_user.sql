-- Add column "display_name" to table: "users"
ALTER TABLE `users` ADD COLUMN `display_name` text NULL;
-- Add column "bio" to table: "users"
ALTER TABLE `users` ADD COLUMN `bio` text NULL;
-- Add column "country" to table: "users"
ALTER TABLE `users` ADD COLUMN `country` text NULL;
-- Add column "social_links" to table: "users"
ALTER TABLE `users` ADD COLUMN `social_links` text NULL;
-- Add column "avatar_url" to table: "users"
ALTER TABLE `users` ADD COLUMN `avatar_url` text NULL;
