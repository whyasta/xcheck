-- xcheck.events definition

CREATE TABLE `events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  `status` int NOT NULL,
  `last_synced_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.imports definition

CREATE TABLE `imports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `upload_file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `imported_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` enum('PENDING','PROCESSING','COMPLETED','ASSIGNED','FAILED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `event_id` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `success_count` bigint unsigned NOT NULL DEFAULT '0',
  `failed_count` bigint unsigned NOT NULL DEFAULT '0',
  `duplicate_count` bigint unsigned NOT NULL DEFAULT '0',
  `failed_values` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `duplicate_values` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.user_roles definition

CREATE TABLE `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `role_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.gates definition

CREATE TABLE `gates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `gate_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `event_id` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `event_id` (`event_id`),
  CONSTRAINT `gates_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.raw_barcodes definition

CREATE TABLE `raw_barcodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `source` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'import, order dll',
  `import_id` bigint unsigned DEFAULT NULL,
  `barcode` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `assign_status` tinyint DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `import_barcodes_imports_FK` (`import_id`),
  CONSTRAINT `import_barcodes_imports_FK` FOREIGN KEY (`import_id`) REFERENCES `imports` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.sessions definition

CREATE TABLE `sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event_id` bigint unsigned DEFAULT NULL,
  `session_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `session_start` datetime NOT NULL,
  `session_end` datetime NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `sessions_event_id_IDX` (`event_id`,`session_start`) USING BTREE,
  CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.ticket_types definition

CREATE TABLE `ticket_types` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_type_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `event_id` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `event_id` (`event_id`),
  CONSTRAINT `ticket_types_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.users definition

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `auth_uuids` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  KEY `users_user_roles_FK` (`role_id`),
  CONSTRAINT `users_user_roles_FK` FOREIGN KEY (`role_id`) REFERENCES `user_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.barcode_logs definition

CREATE TABLE `barcode_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event_id` bigint unsigned DEFAULT NULL,
  `barcode` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ticket_type_id` bigint unsigned NOT NULL,
  `gate_id` bigint unsigned NOT NULL,
  `session_id` bigint unsigned NOT NULL,
  `scanned_at` timestamp NOT NULL,
  `scanned_by` bigint unsigned NOT NULL,
  `action` enum('IN','OUT') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `device` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `barcode_logs_unique` (`event_id`,`barcode`,`scanned_at`,`action`),
  KEY `barcode_logs_barcode_IDX` (`barcode`) USING BTREE,
  KEY `barcode_logs_users_FK` (`scanned_by`),
  KEY `barcode_logs_gate_id_IDX` (`gate_id`) USING BTREE,
  KEY `barcode_logs_ticket_type_id_IDX` (`ticket_type_id`) USING BTREE,
  CONSTRAINT `barcode_logs_events_FK` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`),
  CONSTRAINT `barcode_logs_users_FK` FOREIGN KEY (`scanned_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.barcodes definition

CREATE TABLE `barcodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `barcode` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `event_id` bigint unsigned NOT NULL,
  `ticket_type_id` bigint unsigned NOT NULL,
  `flag` enum('VALID','USED','EXPIRED') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `current_status` enum('','IN','OUT') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `barcodes_event_id_IDX` (`event_id`,`barcode`) USING BTREE,
  KEY `barcodes_ticket_types_FK` (`ticket_type_id`),
  CONSTRAINT `barcodes_ticket_types_FK` FOREIGN KEY (`ticket_type_id`) REFERENCES `ticket_types` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.gate_users definition

CREATE TABLE `gate_users` (
  `gate_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  UNIQUE KEY `gate_users_gate_id_IDX` (`gate_id`,`user_id`) USING BTREE,
  KEY `gate_users_users_FK` (`user_id`),
  CONSTRAINT `gate_users_gates_FK` FOREIGN KEY (`gate_id`) REFERENCES `gates` (`id`),
  CONSTRAINT `gate_users_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.barcode_gates definition

CREATE TABLE `barcode_gates` (
  `barcode_id` bigint unsigned NOT NULL,
  `gate_id` bigint unsigned NOT NULL,
  UNIQUE KEY `barcode_gates_gate_id_IDX` (`gate_id`,`barcode_id`) USING BTREE,
  KEY `barcode_gates_barcodes_FK` (`barcode_id`),
  CONSTRAINT `barcode_gates_barcodes_FK` FOREIGN KEY (`barcode_id`) REFERENCES `barcodes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `barcode_gates_gates_FK` FOREIGN KEY (`gate_id`) REFERENCES `gates` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- xcheck.barcode_sessions definition

CREATE TABLE `barcode_sessions` (
  `barcode_id` bigint unsigned NOT NULL,
  `session_id` bigint unsigned NOT NULL,
  UNIQUE KEY `barcode_sessions_barcode_id_IDX` (`barcode_id`,`session_id`) USING BTREE,
  KEY `barcode_sessions_sessions_FK` (`session_id`),
  CONSTRAINT `barcode_sessions_barcodes_FK` FOREIGN KEY (`barcode_id`) REFERENCES `barcodes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `barcode_sessions_sessions_FK` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO user_roles
(id, role_name, description, created_at, updated_at, deleted_at)
VALUES(1, 'ADM', 'Admin', '2024-07-07 18:29:25', '2024-07-07 19:02:41', NULL);
INSERT INTO user_roles
(id, role_name, description, created_at, updated_at, deleted_at)
VALUES(2, 'CHECKER', 'Checker', '2024-07-08 21:13:36', '2024-07-08 21:13:36', NULL);

INSERT INTO users
(id, username, password, email, role_id, auth_uuids, created_at, updated_at, deleted_at)
VALUES(1, 'admin', '$2a$10$cNEeqEoDvxlgYFuwvqpSneLEV.xGgc1ghKw4bNS2H3eT4trGjwZ9S', 'info@bigmind.id', 1, '["59dfd69f-d6f7-485a-8d5f-1826c678f4cd", "9a63e75a-e9e0-4480-9d04-306bd246abb9", "2b0a5c08-97f5-4cac-bc48-995cca7748eb", "512d8b46-d62a-4430-b737-0bc74b7f72c8", "95f7785e-cb6a-4bb3-ae20-f1cf26809ff4", "5291f64d-33d1-425d-a86b-82a44a45cd37", "7abf75c1-5f87-48e6-900a-af916d856f9a", "3107056f-1f63-4965-a0d1-dac9d1a5103c", "676f22f1-3070-4130-b6c0-10fdb6f25786", "5a93d1cd-7f53-499a-aae2-e8517c041157", "d372ed0a-bf62-4072-9fa6-191f0b9ff60a", "04e429e2-39d4-4223-9aed-e6b667c6a8c2"]', '2024-07-08 22:03:29', '2024-10-04 23:13:50', NULL);

