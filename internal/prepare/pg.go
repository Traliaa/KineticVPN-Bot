package prepare

import (
	"context"
	"fmt"

	"github.com/Traliaa/KineticVPN-Bot/internal/config"
	"github.com/Traliaa/KineticVPN-Bot/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MustNewPg ...
func MustNewPg(ctx context.Context, cfg *config.Config) (*db.PgTxManager, *pgxpool.Pool, error) {
	poolMaster, err := db.NewPool(ctx, db.PoolConfig{
		DSN: cfg.DB,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create poolMaster: %w", err)
	}

	err = poolMaster.Ping(ctx)
	if err != nil {
		return nil, nil, err
	}
	return db.NewPgTxManager(poolMaster), poolMaster, nil
}
