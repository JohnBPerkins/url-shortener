//go:build integration
// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/JohnBPerkins/url-shortener/gen"
	"github.com/JohnBPerkins/url-shortener/modules/flake"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
    svc       gen.ShortenerServer
    ctx       = context.Background()
    testURL   = "example.com/foo"
)

func TestMain(m *testing.M) {
    dsn := os.Getenv("DATABASE_DSN")
    if dsn == "" {
        fmt.Fprintln(os.Stderr, "DATABASE_DSN is required for integration tests")
        os.Exit(1)
    }
    redisAddr := os.Getenv("REDIS_ENDPOINT")
    if redisAddr == "" {
        fmt.Fprintln(os.Stderr, "REDIS_ENDPOINT is required for integration tests")
        os.Exit(1)
    }

    pgPool, err := pgxpool.Connect(ctx, dsn)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to connect to Postgres: %v\n", err)
        os.Exit(1)
    }
    _, err = pgPool.Exec(ctx, `
      CREATE TABLE IF NOT EXISTS links (
        code TEXT PRIMARY KEY,
        url TEXT      NOT NULL,
        created_at TIMESTAMPTZ NOT NULL
      );
      TRUNCATE TABLE links;
    `)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to prepare DB: %v\n", err)
        os.Exit(1)
    }

    // 3) Connect Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    defer rdb.Close()
    // flush any old data
    if err := rdb.FlushDB(ctx).Err(); err != nil {
        fmt.Fprintf(os.Stderr, "failed to flush Redis: %v\n", err)
        os.Exit(1)
    }

    // 4) Init Sonyflake
    sf := flake.NewSonyflake()

    // 5) Build the service under test
    svc = NewShortenerService(pgPool, rdb, sf)

    // 6) Run tests
    code := m.Run()

    os.Exit(code)
}

func TestIntegration_ShrinkThenResolve(t *testing.T) {
    // call Shrink
    resp, err := svc.Shrink(ctx, &gen.ShortenRequest{Url: testURL})
    if err != nil {
        t.Fatalf("Shrink failed: %v", err)
    }
    if len(resp.Code) != codeLength {
        t.Errorf("expected code length %d, got %d", codeLength, len(resp.Code))
    }

    // call Resolve
    res2, err := svc.Resolve(ctx, &gen.ResolveRequest{Code: resp.Code})
    if err != nil {
        t.Fatalf("Resolve failed: %v", err)
    }
    if res2.Url != testURL {
        t.Errorf("expected URL %q, got %q", testURL, res2.Url)
    }
}

func TestIntegration_CacheHit(t *testing.T) {
    // Shrink once
    resp, err := svc.Shrink(ctx, &gen.ShortenRequest{Url: testURL})
    if err != nil {
        t.Fatalf("Shrink failed: %v", err)
    }
    code := resp.Code

    // First Resolve should populate Redis
    _, err = svc.Resolve(ctx, &gen.ResolveRequest{Code: code})
    if err != nil {
        t.Fatalf("first Resolve failed: %v", err)
    }

    // Directly inspect Redis: it should have a key
    rdb := svc.(*ShortenerService).cache // type-assert to grab the redis.Client
    cached, err := rdb.Get(ctx, code).Result()
    if err != nil {
        t.Fatalf("expected cache to contain %s: %v", code, err)
    }
    if cached != testURL {
        t.Errorf("cache[%s]=%q; want %q", code, cached, testURL)
    }
}
