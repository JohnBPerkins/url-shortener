package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// With 3 Railway instances: 3 × 15 = 45 connections (under 60 limit)
	// With 1 instance: 15 connections is plenty
	config.MaxConns = 15

	return pgxpool.ConnectConfig(ctx, config)
}

type Pool = pgxpool.Pool