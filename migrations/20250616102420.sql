-- Create "invoices" table
CREATE TABLE `invoices` (
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
  `vat` text NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_invoices_deleted_at" to table: "invoices"
CREATE INDEX `idx_invoices_deleted_at` ON `invoices` (`deleted_at`);
-- Create "items" table
CREATE TABLE `items` (
  `id` varchar NULL,
  `invoice_id` varchar NULL,
  `name` text NULL,
  `type` text NULL,
  `quantity` integer NULL,
  `tax` integer NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_invoices_items` FOREIGN KEY (`invoice_id`) REFERENCES `invoices` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
