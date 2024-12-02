package isuports

import (
	"fmt"
	"os/exec"
)

const (
	adminDBSchemaFilePath    = "../sql/admin/10_schema.sql"
	adminDBDataFilePath      = "../sql/admin/90_data.sql"
	adminDBMigrationFilePath = "../sql/admin/99_migration.sql"
)

func executeQueryFile(path string) error {
	script := fmt.Sprintf(
		"mysql -u %s -p %s --host %s --port %s %s < %s",
		getEnv("ISUCON_DB_USER", "isucon"),
		getEnv("ISUCON_DB_PASSWORD", "isucon"),
		getEnv("ISUCON_DB_HOST", "127.0.0.1"),
		getEnv("ISUCON_DB_PORT", "3306"),
		getEnv("ISUCON_DB_NAME", "isuports"),
		path,
	)
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
	return nil
}
