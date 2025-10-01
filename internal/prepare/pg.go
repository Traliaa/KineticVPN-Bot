package prepare

import (
	"context"
	"fmt"

	"github.com/Traliaa/KineticVPN-Bot/internal/config"
	"github.com/Traliaa/KineticVPN-Bot/pkg/db"
)

// MustNewPg ...
func MustNewPg(ctx context.Context, cfg *config.Config) (*db.PgTxManager, error) {
	poolMaster, err := db.NewPool(ctx, db.PoolConfig{
		DSN: cfg.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create poolMaster: %w", err)
	}
	return db.NewPgTxManager(poolMaster), nil
}
