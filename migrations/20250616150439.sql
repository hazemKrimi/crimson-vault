-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_invoices" table
CREATE TABLE `new_invoices` (
  `id` varchar NULL,
  `user_id` varchar NULL,
  `client_id` varchar NULL,
  `created_at` datetime NULL,
  `due_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `reference` text NULL,
  `status` text NULL,
  `pdf` text NULL,
  `currency` text NULL,
  `vat` integer NULL,
  PRIMARY KEY (`id`)
);
-- Copy rows from old table "invoices" to new temporary table "new_invoices"
INSERT INTO `new_invoices` (`id`, `user_id`, `client_id`, `created_at`, `due_at`, `updated_at`, `deleted_at`, `reference`, `status`, `pdf`, `currency`, `vat`) SELECT `id`, `user_id`, `client_id`, `created_at`, `due_at`, `updated_at`, `deleted_at`, `reference`, `status`, `pdf`, `currency`, `vat` FROM `invoices`;
-- Drop "invoices" table after copying rows
DROP TABLE `invoices`;
-- Rename temporary table "new_invoices" to "invoices"
ALTER TABLE `new_invoices` RENAME TO `invoices`;
-- Create index "idx_invoices_deleted_at" to table: "invoices"
CREATE INDEX `idx_invoices_deleted_at` ON `invoices` (`deleted_at`);
-- Add column "user_id" to table: "items"
ALTER TABLE `items` ADD COLUMN `user_id` varchar NULL;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
