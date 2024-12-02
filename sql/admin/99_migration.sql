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

DROP TABLE IF EXISTS `player_rank`;
CREATE TABLE `player_rank` (
	`player_id` VARCHAR(255) NOT NULL,
	`competition_id` VARCHAR(255) NOT NULL,
	`competition_title` TEXT NOT NULL,
	`competition_created_at` BIGINT NOT NULL,
	`score` BIGINT NOT NULL,
	UNIQUE KEY `player_competition` (`player_id`, `competition_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO player_rank (player_id, competition_id, competition_title, competition_created_at, score)
SELECT ps.player_id AS player_id, c.id AS competition_id, c.title AS competition_title, c.created_at AS competition_created_at, ps.score AS score
FROM player_score AS ps
INNER JOIN competition AS c ON c.id = ps.competition_id;

ALTER TABLE player_score
DROP INDEX competition_player_idx,
ADD INDEX competition_player_idx (competition_id, player_id),
DROP INDEX competition_score_row_idx,
ADD INDEX competition_score_row_idx (competition_id, score DESC, row_num ASC);

ALTER TABLE player
DROP INDEX `tenant_created_idx`,
ADD INDEX `tenant_created_idx` (`tenant_id`, `created_at` DESC)
