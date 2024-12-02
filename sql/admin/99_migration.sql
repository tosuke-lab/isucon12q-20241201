USE `isuports`;

DROP TABLE IF EXISTS `first_visit`;

CREATE TABLE `first_visit` (
	`player_id` VARCHAR(255) NOT NULL,
	`tenant_id` BIGINT UNSIGNED NOT NULL,
	`competition_id` VARCHAR(255) NOT NULL,
	`created_at` BIGINT NOT NULL,
	UNIQUE KEY `player_tenant_competition` (`player_id`, `tenant_id`, `competition_id`),
	INDEX `competition_idx` (`competition_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO first_visit(player_id, tenant_id, competition_id, created_at)
SELECT player_id, tenant_id, competition_id, MIN(created_at) AS created_at
FROM visit_history
GROUP BY player_id, tenant_id, competition_id;

ALTER TABLE player_score
DROP INDEX competition_player_idx,
ADD INDEX competition_player_idx (competition_id, player_id),
DROP INDEX competition_score_row_idx,
ADD INDEX competition_score_row_idx (competition_id, score DESC, row_num ASC);

ALTER TABLE player
DROP INDEX `tenant_created_idx`,
ADD INDEX `tenant_created_idx` (`tenant_id`, `created_at` DESC)
