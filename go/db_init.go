package isuports

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
)

const (
	adminDBSchemaFilePath    = "../sql/admin/10_schema.sql"
	adminDBDataFilePath      = "../sql/admin/90_data.sql"
	adminDBMigrationFilePath = "../sql/admin/99_migration.sql"
)

func executeQueryFile(path string) error {
	script := fmt.Sprintf(
		"mysql -u\"%s\" -p\"%s\" --host \"%s\" --port \"%s\" %s < %s",
		getEnv("ISUCON_DB_USER", "isucon"),
		getEnv("ISUCON_DB_PASSWORD", "isucon"),
		getEnv("ISUCON_DB_HOST", "127.0.0.1"),
		getEnv("ISUCON_DB_PORT", "3306"),
		getEnv("ISUCON_DB_NAME", "isuports"),
		path,
	)
	log.Printf("executing query: %s", script)
	cmd := exec.Command("sh", "-c", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to execute query file: %s: %s: %w", path, out, err)
	}
	return nil
}

func CleanDB() error {
	if err := executeQueryFile(adminDBSchemaFilePath); err != nil {
		return err
	}
	if err := executeQueryFile(adminDBDataFilePath); err != nil {
		return err
	}
	if err := executeQueryFile(adminDBMigrationFilePath); err != nil {
		return err
	}

	var err error
	adminDB, err = connectAdminDB()
	if err != nil {
		return err
	}

	if err := migrateAllTenantDB(); err != nil {
		return err
	}

	return nil
}

func migrateAllTenantDB() error {
	var ts []TenantRow
	if err := adminDB.Select(&ts, "SELECT * FROM tenant"); err != nil {
		return fmt.Errorf("failed to select from tenant: %w", err)
	}
	var wg sync.WaitGroup
	for _, t := range ts {
		wg.Add(1)
		go func(t TenantRow) {
			defer wg.Done()
			if err := migrateTenantDB(t.ID); err != nil {
				log.Printf("failed to migrate tenant DB: tenant_id=%d: %s", t.ID, err)
			}
		}(t)
	}
	wg.Wait()
	return nil
}

func migrateTenantDB(id int64) error {
	tenantDB, err := connectToInitialTenantDB(id)
	if err != nil {
		return err
	}

	var cs []CompetitionRow
	if err := tenantDB.Select(&cs, "SELECT * FROM competition"); err != nil {
		return fmt.Errorf("failed to select from competition: %w", err)
	}
	if _, err := adminDB.NamedExec(
		"INSERT INTO competition (tenant_id, id, title, finished_at, created_at, updated_at) VALUES (:tenant_id, :id, :title, :finished_at, :created_at, :updated_at)",
		cs,
	); err != nil {
		return fmt.Errorf("failed to insert into competition: %w", err)
	}

	var ps []PlayerRow
	if err := tenantDB.Select(&ps, "SELECT * FROM player"); err != nil {
		return fmt.Errorf("failed to select from player: %w", err)
	}
	if _, err := adminDB.NamedExec(
		"INSERT INTO player (tenant_id, id, display_name, is_disqualified, created_at, updated_at) VALUES (:tenant_id, :id, :display_name, :is_disqualified, :created_at, :updated_at)",
		ps,
	); err != nil {
		return fmt.Errorf("failed to insert into player: %w", err)
	}

	var pss []PlayerScoreRow
	if err := tenantDB.Select(&pss, "SELECT * FROM player_score ORDER BY row_num DESC"); err != nil {
		return fmt.Errorf("failed to select from player_score: %w", err)
	}
	type key struct {
		playerID      string
		competitionID string
	}
	pssMap := map[key]PlayerScoreRow{}
	for _, ps := range pss {
		k := key{playerID: ps.PlayerID, competitionID: ps.CompetitionID}
		if _, ok := pssMap[k]; ok {
			continue
		}
		pssMap[k] = ps
	}
	fpss := make([]PlayerScoreRow, 0, len(pssMap))
	for _, ps := range pssMap {
		fpss = append(fpss, ps)
	}

	for {
		var psss []PlayerScoreRow
		if len(fpss) > 1000 {
			psss = fpss[:1000]
			fpss = fpss[1000:]
		} else if len(fpss) > 0 {
			psss = fpss
			fpss = nil
		} else {
			break
		}
		if _, err := adminDB.NamedExec(
			"INSERT INTO player_score (tenant_id, id, player_id, competition_id, score, row_num, created_at, updated_at) VALUES (:tenant_id, :id, :player_id, :competition_id, :score, :row_num, :created_at, :updated_at)",
			psss,
		); err != nil {
			return fmt.Errorf("failed to insert into player_score: %w", err)
		}
	}

	log.Printf("migrated tenant DB: tenant_id=%d", id)

	return nil
}

func connectToInitialTenantDB(id int64) (*sqlx.DB, error) {
	p := filepath.Join("../../initial_data", fmt.Sprintf("%d.db", id))
	db, err := sqlx.Open(sqliteDriverName, fmt.Sprintf("file:%s?mode=ro", p))
	if err != nil {
		return nil, fmt.Errorf("failed to open initial tenant DB: %s: %w", p, err)
	}
	return db, nil
}
