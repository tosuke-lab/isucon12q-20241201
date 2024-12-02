package isuports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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

func (p *PlayerRow) isFresh() bool {
	updatedAt := time.Unix(p.UpdatedAt, 0)
	return time.Now().Sub(updatedAt) < 2*time.Second
}

type playerEntry struct {
	sync.Mutex
	r *PlayerRow
}

func (p *playerEntry) retrieve(ctx context.Context, id string) (*PlayerRow, error) {
	p.Lock()
	defer p.Unlock()

	if p.r == nil || !p.r.isFresh() {
		pe, err := forceRetrievePlayer(ctx, id)
		if err != nil {
			p.r = nil
			return nil, err
		}
		p.r = pe
	}
	return p.r, nil
}

var (
	playerMap   = map[string]*playerEntry{}
	playerMapMu = sync.RWMutex{}
)

// 参加者を取得する
func retrievePlayer(ctx context.Context, id string) (*PlayerRow, error) {
	playerMapMu.RLock()
	v, ok := playerMap[id]
	playerMapMu.RUnlock()
	if ok {
		return v.retrieve(ctx, id)
	}

	playerMapMu.Lock()
	v, ok = playerMap[id]
	if !ok {
		v = &playerEntry{}
		playerMap[id] = v
	}
	playerMapMu.Unlock()
	return v.retrieve(ctx, id)
}

func forceRetrievePlayer(ctx context.Context, id string) (*PlayerRow, error) {
	var p PlayerRow
	if err := adminDB.GetContext(ctx, &p, "SELECT * FROM player WHERE id = ?", id); err != nil {
		return nil, fmt.Errorf("error Select player: id=%s, %w", id, err)
	}
	return &p, nil
}

// 参加者を認可する
// 参加者向けAPIで呼ばれる
func authorizePlayer(ctx context.Context, id string) error {
	player, err := retrievePlayer(ctx, id)
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
func retrieveCompetition(ctx context.Context, tenantID int64, id string) (*CompetitionRow, error) {
	var c CompetitionRow
	if err := adminDB.GetContext(ctx, &c, "SELECT * FROM competition WHERE tenant_id = ? AND id = ?", tenantID, id); err != nil {
		return nil, fmt.Errorf("error Select competition: id=%s, %w", id, err)
	}
	return &c, nil
}

// tenant/player_score
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
