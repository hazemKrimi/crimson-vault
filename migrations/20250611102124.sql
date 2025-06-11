-- Create "clients" table
CREATE TABLE `clients` (
  `id` varchar NULL,
  `user_id` varchar NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NULL,
  `fiscal_code` text NULL,
  `address` text NULL,
  `zip` text NULL,
  `country` text NULL,
  `phone` text NULL,
  `email` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_clients` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_clients_deleted_at" to table: "clients"
CREATE INDEX `idx_clients_deleted_at` ON `clients` (`deleted_at`);
-- Create "users" table
CREATE TABLE `users` (
  `id` varchar NULL,
  `session_id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `logo` text NULL,
  `name` text NULL,
  `fiscal_code` text NULL,
  `address` text NULL,
  `zip` text NULL,
  `country` text NULL,
  `phone` text NULL,
  `email` text NULL,
  `username` text NULL,
  `password` text NULL,
  PRIMARY KEY (`id`)
);
-- Create index "users_username" to table: "users"
CREATE UNIQUE INDEX `users_username` ON `users` (`username`);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
