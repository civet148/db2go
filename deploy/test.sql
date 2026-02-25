CREATE DATABASE `test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `test`;


CREATE TABLE `inventory_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_frozen` tinyint(1) DEFAULT '0',
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品：名称；不能为空',
  `serial_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品序列号',
  `quantity` decimal(16,3) DEFAULT '0.000',
  `price` decimal(16,2) DEFAULT '0.00',
  `location` point DEFAULT NULL,
  `product_extra` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventory_data_created_at` (`created_at`),
  KEY `idx_inventory_data_updated_at` (`updated_at`),
  KEY `i_serial_no` (`serial_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `inventory_in` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_deleted` tinyint(1) DEFAULT '0',
  `delete_time` datetime DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `order_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT '0',
  `user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `quantity` decimal(16,6) DEFAULT '0.000000',
  `weight` decimal(16,6) DEFAULT '0.000000',
  `remark` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`),
  KEY `idx_inventory_in_created_at` (`created_at`),
  KEY `idx_inventory_in_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `inventory_out` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_deleted` tinyint(1) DEFAULT '0',
  `delete_time` datetime DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT '0',
  `order_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT '0',
  `user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `quantity` decimal(16,6) DEFAULT '0.000000',
  `weight` decimal(16,6) DEFAULT '0.000000',
  `remark` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_PROD_USER` (`product_id`,`user_id`),
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`),
  KEY `idx_inventory_out_created_at` (`created_at`),
  KEY `idx_inventory_out_updated_at` (`updated_at`),
  KEY `i_product_id` (`product_id`),
  KEY `i_user_id` (`user_id`),
  KEY `FULTXT_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`),
  KEY `idx_roles_created_at` (`created_at`),
  KEY `idx_roles_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `user_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` bigint unsigned DEFAULT NULL,
  `avatar` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `address` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_profiles_user_id` (`user_id`),
  KEY `idx_user_profiles_created_at` (`created_at`),
  KEY `idx_user_profiles_updated_at` (`updated_at`),
  CONSTRAINT `fk_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_user_roles_role` (`role_id`),
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_name` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `email` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_user_name` (`user_name`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_created_at` (`created_at`),
  KEY `idx_users_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
