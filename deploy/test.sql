CREATE DATABASE `test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `test`;


CREATE TABLE `inventory_data` (
  `id` bigint unsigned NOT NULL COMMENT '产品ID',
  `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
  `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_frozen` tinyint(1) NOT NULL DEFAULT '0' COMMENT '冻结状态(0: 未冻结 1: 已冻结)',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '产品名称',
  `serial_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '产品编号',
  `quantity` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '产品库存',
  `price` decimal(16,2) NOT NULL DEFAULT '0.00' COMMENT '产品均价',
  `product_extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '产品附带数据(JSON文本)',
  `location` point DEFAULT NULL COMMENT '地理位置',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品库存数据表';

CREATE TABLE `inventory_in` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
  `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除状态(0: 未删除 1: 已删除)',
  `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
  `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '入库单号',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '交货人ID',
  `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交货人姓名',
  `quantity` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '数量',
  `weight` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '净重',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='入库主表';

CREATE TABLE `inventory_out` (
  `id` bigint unsigned NOT NULL COMMENT '主键ID',
  `create_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '创建人姓名',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新人ID',
  `update_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '更新人姓名',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除状态(0: 未删除 1: 已删除)',
  `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
  `product_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '产品ID',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '出库单号',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收货人ID',
  `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `quantity` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '数量',
  `weight` decimal(16,6) NOT NULL DEFAULT '0.000000' COMMENT '净重',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='出库主表';
