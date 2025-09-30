package db

import (
	"context"
	"fmt"

	"github.com/Traliaa/KineticVPN-Bot/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dsnFormat = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

type PoolConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type PgTxManager struct {
	poolMaster  *pgxpool.Pool
	poolReplica *pgxpool.Pool
}

func NewPgTxManager(poolMaster, poolReplica *pgxpool.Pool) *PgTxManager {
	return &PgTxManager{
		poolMaster:  poolMaster,
		poolReplica: poolReplica,
	}
}

func (m *PgTxManager) Close() {
	m.poolReplica.Close()
	m.poolMaster.Close()
}

func NewPool(ctx context.Context, conf PoolConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(dsnFormat, conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	return pgxpool.New(ctx, dsn)
}

func (m *PgTxManager) RunMaster(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error {
	options := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	// то что запрос нужно выполнить на мастере еще не означает что это нужно выполнить в транзакции, может требоваться
	// просто согласованное чтение, например.
	return m.inTx(ctx, m.poolMaster, options, fn)
}

func (m *PgTxManager) RunReplica(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error {
	options := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	return m.inTx(ctx, m.poolReplica, options, fn)
}

//func (m *PgTxManager) Conn() Transaction {
//	return m.poolMaster
//}

func (m *PgTxManager) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error {
	options := pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	}
	return m.inTx(ctx, m.poolReplica, options, fn)
}

func (m *PgTxManager) inTx(
	ctx context.Context,
	pool *pgxpool.Pool,
	options pgx.TxOptions,
	f func(ctxTx context.Context, tx Transaction) error,
) error {
	tx, err := pool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to begin tx, err: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Info("%v", p)
			_ = tx.Rollback(ctx)
			panic(p) // fallthrough panic after rollback on caught panic
		} else if err != nil {
			_ = tx.Rollback(ctx) // if error during computations
		} else {
			err = tx.Commit(ctx) // all good
		}
	}()

	err = f(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to run fn, err: %w", err)
	}

	return nil
}
