-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_items" table
CREATE TABLE `new_items` (
  `id` varchar NULL,
  `invoice_id` varchar NULL,
  `user_id` varchar NULL,
  `name` text NULL,
  `type` text NULL,
  `price` real NULL,
  `quantity` integer NULL,
  `tax` integer NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_invoices_items` FOREIGN KEY (`invoice_id`) REFERENCES `invoices` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "items" to new temporary table "new_items"
INSERT INTO `new_items` (`id`, `invoice_id`, `user_id`, `name`, `type`, `price`, `quantity`, `tax`) SELECT `id`, `invoice_id`, `user_id`, `name`, `type`, `price`, `quantity`, `tax` FROM `items`;
-- Drop "items" table after copying rows
DROP TABLE `items`;
-- Rename temporary table "new_items" to "items"
ALTER TABLE `new_items` RENAME TO `items`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
