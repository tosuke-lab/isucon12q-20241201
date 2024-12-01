USE `isuports`;

DROP TABLE IF EXISTS `first_visit`;
DROP TABLE IF EXISTS `billing_report`;

CREATE TABLE `first_visit` (
  `player_id` VARCHAR(255) NOT NULL,
  `tenant_id` BIGINT UNSIGNED NOT NULL,
  `competition_id` VARCHAR(255) NOT NULL,
  `created_at` BIGINT NOT NULL,
  INDEX `tenant_id_idx` (`tenant_id`),
  UNIQUE KEY `player_tenant_competition` (`player_id`, `tenant_id`, `competition_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO first_visit(player_id, tenant_id, competition_id, created_at)
SELECT player_id, tenant_id, competition_id, MIN(created_at) AS created_at
FROM visit_history
GROUP BY player_id, tenant_id, competition_id;

CREATE TABLE `billing_report` (
	`tenant_id` BIGINT UNSIGNED NOT NULL,
	`competition_id` VARCHAR(255) NOT NULL,
	`competition_title` TEXT NOT NULL,
	`player_count` BIGINT NOT NULL,
	`visitor_count` BIGINT NOT NULL,
	`billing_player_yen` BIGINT NOT NULL,
	`billing_visitor_yen` BIGINT NOT NULL,
	`billing_yen` BIGINT NOT NULL,
	INDEX `tenant_id_idx` (`tenant_id`),
	UNIQUE KEY `tenant_competition` (`tenant_id`, `competition_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;
