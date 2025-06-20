package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
  return pgxpool.Connect(ctx, dsn)
}

type Pool = pgxpool.Pool