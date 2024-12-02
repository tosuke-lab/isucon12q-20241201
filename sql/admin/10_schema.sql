USE `isuports`;

DROP TABLE IF EXISTS `tenant`;
DROP TABLE IF EXISTS `id_generator`;
DROP TABLE IF EXISTS `visit_history`;
DROP TABLE IF EXISTS `competition`;
DROP TABLE IF EXISTS `player`;
DROP TABLE IF EXISTS `player_score`;
DROP TABLE IF EXISTS `billing_report`;

CREATE TABLE `tenant` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `display_name` VARCHAR(255) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `id_generator` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `stub` CHAR(1) NOT NULL DEFAULT '',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `stub` (`stub`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `visit_history` (
  `player_id` VARCHAR(255) NOT NULL,
  `tenant_id` BIGINT UNSIGNED NOT NULL,
  `competition_id` VARCHAR(255) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  INDEX `tenant_id_idx` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `competition` (
	`id` VARCHAR(255) NOT NULL PRIMARY KEY,
	`tenant_id` BIGINT NOT NULL,
	`title` TEXT NOT NULL,
	`finished_at` BIGINT NULL,
	`created_at` BIGINT NOT NULL,
	`updated_at` BIGINT NOT NULL,
	INDEX `tenant_id_idx` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `player` (
	`id` VARCHAR(255) NOT NULL PRIMARY KEY,
	`tenant_id` BIGINT NOT NULL,
	`display_name` TEXT NOT NULL,
	`is_disqualified` BOOLEAN NOT NULL,
	`created_at` BIGINT NOT NULL,
	`updated_at` BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `player_score` (
	`id` VARCHAR(255) NOT NULL PRIMARY KEY,
	`tenant_id` BIGINT NOT NULL,
	`player_id` VARCHAR(255) NOT NULL,
	`competition_id` VARCHAR(255) NOT NULL,
	`score` BIGINT NOT NULL,
	`row_num` BIGINT NOT NULL,
	`created_at` BIGINT NOT NULL,
	`updated_at` BIGINT NOT NULL,
	INDEX `tenant_id_idx` (`tenant_id`),
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `billing_report` (
	`tenant_id` BIGINT UNSIGNED NOT NULL,
	`competition_id` VARCHAR(255) NOT NULL,
	`competition_title` TEXT NOT NULL,
	`player_count` BIGINT NOT NULL,
	`visitor_count` BIGINT NOT NULL,
	`billing_player_yen` BIGINT NOT NULL,
	`billing_visitor_yen` BIGINT NOT NULL,
	`billing_yen` BIGINT NOT NULL,
	`created_at` BIGINT NOT NULL,
	INDEX `tenant_id_idx` (`tenant_id`),
	UNIQUE KEY `tenant_competition` (`tenant_id`, `competition_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;
