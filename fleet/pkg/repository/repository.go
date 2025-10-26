package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repository struct {
	*Queries
	snowflakeNode *snowflake.Node
	pgxPool       *pgxpool.Pool
	redisConn     *redis.Client
}

type Tx struct {
	Tx      pgx.Tx
	Queries *Queries
}

func CreateRepository(pool *pgxpool.Pool, redisConn *redis.Client, logger *zap.Logger) (repo *Repository, err error) {
	snowflakeNode, err := snowflake.NewNode(1)
	if err != nil {
		logger.Fatal("Failed to create snowflake node", zap.Error(err))
	}

	return &Repository{
		Queries:       New(pool),
		snowflakeNode: snowflakeNode,
		pgxPool:       pool,
		redisConn:     redisConn,
	}, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pgxPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *Repository) WithTx(tx pgx.Tx) Querier {
	return r.Queries.WithTx(tx)
}

func (r *Repository) CommitTx(tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

func (r *Repository) RollbackTx(tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	if err := tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
		return err
	}
	return nil
}

func (r *Repository) GenerateSnowflakeID() snowflake.ID {
	return r.snowflakeNode.Generate()
}
