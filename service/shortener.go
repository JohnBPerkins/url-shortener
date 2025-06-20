package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/JohnBPerkins/url-shortener/gen"
	"github.com/JohnBPerkins/url-shortener/modules/db"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sony/sonyflake"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)	

const (
	maxURLLength = 2048
	maxAttempts = 5
	codeLength = 8
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type ShortenerService struct {
	gen.UnimplementedShortenerServer
	dbPool *db.Pool
	cache *redis.Client
	flake *sonyflake.Sonyflake
}

func newShortenerService(dbPool *db.Pool, cache *redis.Client ,flake *sonyflake.Sonyflake) gen.ShortenerServer {
	return &ShortenerService{dbPool: dbPool, cache: cache, flake: flake}
}

func (s *ShortenerService) Shrink(ctx context.Context, req *gen.ShrinkRequest) (*gen.ShrinkResponse, error) {
	if !isValidURL(req.GetUrl()) {
        return nil, status.Errorf(codes.InvalidArgument, "invalid URL: %q", req.GetUrl())
    }

	for i := 0; i < maxAttempts; i++ {
		id, err := s.flake.NextID()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate ID: %v", err)
		}
		code := encodeBase62(id)

		_, err = s.dbPool.Exec(ctx,
			`INSERT INTO links (code, url, created_at) VALUES ($1, $2, NOW())`,
			code, req.GetUrl(),
		)
		if err == nil {
			if err := s.cache.Set(ctx, code, req.GetUrl(), 24*time.Hour).Err(); err != nil {}
			return &gen.ShrinkResponse{Code: code}, nil
		}
		if isUniqueViolation(err) {
            continue
        }
		return nil, status.Errorf(codes.Internal, "db insert failed: %v", err)
	}

	return nil, status.Errorf(codes.Internal,
        "could not generate a unique code after %d attempts", maxAttempts)
}

func isValidURL(candidate string) bool {
    if len(candidate) == 0 || len(candidate) > maxURLLength {
        return false
    }

    raw := candidate
    if !strings.HasPrefix(candidate, "http://") && !strings.HasPrefix(candidate, "https://") {
        raw = "http://" + candidate
    }

    u, err := url.ParseRequestURI(raw)
    if err != nil {
        return false
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return false
    }
    if u.Host == "" {
        return false
    }
    return true
}

func isUniqueViolation(err error) bool {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        return true
    }
    return false
}

func encodeBase62(num uint64) string {
    var encoded []byte
    if num == 0 {
        encoded = []byte{base62Chars[0]}
    } else {
        for num > 0 {
            rem := num % 62
            encoded = append(encoded, base62Chars[rem])
            num /= 62
        }
        for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
            encoded[i], encoded[j] = encoded[j], encoded[i]
        }
    }

    if len(encoded) < codeLength {
        padding := make([]byte, codeLength-len(encoded))
        for i := range padding {
            padding[i] = base62Chars[0]
        }
        encoded = append(padding, encoded...)
    } else if len(encoded) > codeLength {
        encoded = encoded[len(encoded)-codeLength:]
    }

    return string(encoded)
}

func (s *ShortenerService) Resolve(ctx context.Context, req *gen.ResolveRequest) (*gen.ResolveResponse, error) {
	code := req.GetCode()
	
	urlStr, err := s.cache.Get(ctx, code).Result()
    if err == nil {
        return &gen.ResolveResponse{Url: urlStr}, nil
    }
    if err != redis.Nil {
        return nil, status.Errorf(codes.Internal, "cache lookup failed: %v", err)
    }

	var dbURL string
    err = s.dbPool.QueryRow(ctx,
        `SELECT url FROM links WHERE code = $1`, code,
    ).Scan(&dbURL)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, status.Errorf(codes.NotFound, "code not found: %s", code)
        }
        return nil, status.Errorf(codes.Internal, "db query failed: %v", err)
    }

	if err := s.cache.Set(ctx, code, dbURL, 24*time.Hour).Err(); err != nil {}
	
    return &gen.ResolveResponse{Url: dbURL}, nil
}
