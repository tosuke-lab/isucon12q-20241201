DELETE FROM tenant WHERE id > 100;
DELETE FROM visit_history WHERE created_at >= '1654041600';
DELETE FROM first_visit WHERE created_at >= '1654041600';
DELETE FROM billing_report WHERE created_at >= '1654041600';
DELETE FROM competition WHERE created_at >= '1654041600';
DELETE FROM player WHERE created_at >= '1654041600';
DELETE FROM player_score WHERE created_at >= '1654041600';
UPDATE id_generator SET id=2678400000 WHERE stub='a';
ALTER TABLE id_generator AUTO_INCREMENT=2678400000;
