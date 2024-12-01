package isuports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// admin/tenant
type TenantRow struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	DisplayName string `db:"display_name"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

// 参加者
// tenant/player
type PlayerRow struct {
	TenantID       int64  `db:"tenant_id"`
	ID             string `db:"id"`
	DisplayName    string `db:"display_name"`
	IsDisqualified bool   `db:"is_disqualified"`
	CreatedAt      int64  `db:"created_at"`
	UpdatedAt      int64  `db:"updated_at"`
}

// 参加者を取得する
func retrievePlayer(ctx context.Context, tenantDB dbOrTx, id string) (*PlayerRow, error) {
	var p PlayerRow
	if err := tenantDB.GetContext(ctx, &p, "SELECT * FROM player WHERE id = ?", id); err != nil {
		return nil, fmt.Errorf("error Select player: id=%s, %w", id, err)
	}
	return &p, nil
}

// 参加者を認可する
// 参加者向けAPIで呼ばれる
func authorizePlayer(ctx context.Context, tenantDB dbOrTx, id string) error {
	player, err := retrievePlayer(ctx, tenantDB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, "player not found")
		}
		return fmt.Errorf("error retrievePlayer from viewer: %w", err)
	}
	if player.IsDisqualified {
		return echo.NewHTTPError(http.StatusForbidden, "player is disqualified")
	}
	return nil
}

// 大会
// tenant/competition
type CompetitionRow struct {
	TenantID   int64         `db:"tenant_id"`
	ID         string        `db:"id"`
	Title      string        `db:"title"`
	FinishedAt sql.NullInt64 `db:"finished_at"`
	CreatedAt  int64         `db:"created_at"`
	UpdatedAt  int64         `db:"updated_at"`
}

// 大会を取得する
func retrieveCompetition(ctx context.Context, tenantDB dbOrTx, id string) (*CompetitionRow, error) {
	var c CompetitionRow
	if err := tenantDB.GetContext(ctx, &c, "SELECT * FROM competition WHERE id = ?", id); err != nil {
		return nil, fmt.Errorf("error Select competition: id=%s, %w", id, err)
	}
	return &c, nil
}

type PlayerScoreRow struct {
	TenantID      int64  `db:"tenant_id"`
	ID            string `db:"id"`
	PlayerID      string `db:"player_id"`
	CompetitionID string `db:"competition_id"`
	Score         int64  `db:"score"`
	RowNum        int64  `db:"row_num"`
	CreatedAt     int64  `db:"created_at"`
	UpdatedAt     int64  `db:"updated_at"`
}
