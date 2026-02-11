CREATE DATABASE `test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `test`;


CREATE TABLE `inventory_data` (
  `id` bigint unsigned NOT NULL COMMENT '产品ID',
  `create_id` bigint unsigned DEFAULT NULL,
  `create_name` longtext COLLATE utf8mb4_unicode_ci,
  `create_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `update_id` bigint unsigned DEFAULT NULL,
  `update_name` longtext COLLATE utf8mb4_unicode_ci,
  `update_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `is_frozen` tinyint(1) DEFAULT '0',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '产品：名称；不能为空',
  `serial_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '产品序列号',
  `quantity` decimal(16,3) DEFAULT '0.000',
  `price` decimal(16,2) DEFAULT '0.00',
  `location` point DEFAULT NULL,
  `product_extra` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `i_serial_no` (`serial_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品库存数据表';

CREATE TABLE `inventory_in` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `create_id` bigint unsigned DEFAULT NULL,
  `create_name` longtext COLLATE utf8mb4_unicode_ci,
  `create_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `update_id` bigint unsigned DEFAULT NULL,
  `update_name` longtext COLLATE utf8mb4_unicode_ci,
  `update_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `is_deleted` tinyint(1) DEFAULT '0',
  `delete_time` datetime DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT '0',
  `user_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `quantity` decimal(16,6) DEFAULT '0.000000',
  `weight` decimal(16,6) DEFAULT '0.000000',
  `remark` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='入库主表';

CREATE TABLE `inventory_out` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `create_id` bigint unsigned DEFAULT NULL,
  `create_name` longtext COLLATE utf8mb4_unicode_ci,
  `create_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `update_id` bigint unsigned DEFAULT NULL,
  `update_name` longtext COLLATE utf8mb4_unicode_ci,
  `update_time` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'CURRENT_TIMESTAMP',
  `is_deleted` tinyint(1) DEFAULT '0',
  `delete_time` datetime DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT '0',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT '0',
  `user_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `quantity` decimal(16,6) DEFAULT '0.000000',
  `weight` decimal(16,6) DEFAULT '0.000000',
  `remark` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE,
  UNIQUE KEY `UNIQ_PROD_USER` (`product_id`,`user_id`,`update_time`),
  KEY `i_user_id` (`user_id`) USING BTREE,
  KEY `i_product_id` (`product_id`) USING BTREE,
  FULLTEXT KEY `FULTXT_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='出库主表';

CREATE TABLE `user_role_id` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_role_id_unique` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
