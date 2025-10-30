package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	pb "github.com/JohnBPerkins/url-shortener/gen"
	"github.com/JohnBPerkins/url-shortener/internal/service"
	"github.com/JohnBPerkins/url-shortener/internal/web"
	"github.com/JohnBPerkins/url-shortener/modules/db"
	"github.com/JohnBPerkins/url-shortener/modules/flake"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {    
	prometheus.MustRegister(
        service.ResolveHits,
        service.ResolveMisses, 
        service.ResolveErrors,
        service.ResolveDuration,
    )

	//init db
	ctx := context.Background()

	// Log available env vars for debugging
	log.Printf("DATABASE_URL: %q", os.Getenv("DATABASE_URL"))
	log.Printf("DATABASE_DSN: %q", os.Getenv("DATABASE_DSN"))
	log.Printf("PGDATABASE: %q", os.Getenv("PGDATABASE"))
	log.Printf("DATABASE_PRIVATE_URL: %q", os.Getenv("DATABASE_PRIVATE_URL"))

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_PRIVATE_URL") // Railway private URL
	}
	if dsn == "" {
		dsn = os.Getenv("DATABASE_DSN") // fallback
	}
	if dsn == "" {
		log.Fatal("No database connection string found in environment variables")
	}
	dbPool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to Postgres (%q): %v", dsn, err)
	}

	//init cache
	log.Printf("REDIS_URL: %q", os.Getenv("REDIS_URL"))
	log.Printf("REDIS_PRIVATE_URL: %q", os.Getenv("REDIS_PRIVATE_URL"))
	log.Printf("REDIS_ENDPOINT: %q", os.Getenv("REDIS_ENDPOINT"))

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = os.Getenv("REDIS_PRIVATE_URL")
	}

	var redisAddr, redisPassword string
	if redisURL != "" {
		// Parse Railway Redis URL: redis://default:password@host:port
		if strings.HasPrefix(redisURL, "redis://") {
			redisURL = strings.TrimPrefix(redisURL, "redis://")
			if strings.Contains(redisURL, "@") {
				parts := strings.Split(redisURL, "@")
				if len(parts) == 2 {
					credParts := strings.Split(parts[0], ":")
					if len(credParts) == 2 {
						redisPassword = credParts[1]
					}
					redisAddr = parts[1]
				}
			} else {
				redisAddr = redisURL
			}
		}
	} else {
		redisAddr = os.Getenv("REDIS_ENDPOINT")
		redisPassword = os.Getenv("REDIS_PASSWORD")
	}

	if redisAddr == "" {
		log.Fatal("No Redis connection string found in environment variables")
	}

	cache := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       0,
    })

	flake := flake.NewSonyflake()
	svc := service.NewShortenerService(dbPool, cache, flake)

	gRpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prom.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prom.StreamServerInterceptor),
	)

	shrinkHandler := web.NewShrinkHandler(svc)
	resolveHandler := web.NewResolveHandler(svc)

	pb.RegisterShortenerServer(gRpcServer, svc)
	grpc_prom.EnableHandlingTimeHistogram()

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Add CORS middleware
	corsHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	mux.HandleFunc("/api/shorten", corsHandler(shrinkHandler))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", resolveHandler)


	go func() {
        log.Println("▶ HTTP API listening on :8080")
        log.Fatal(http.ListenAndServe(":8080", mux))
    }()

	lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Printf("gRPC server listening on %s", lis.Addr())
    log.Fatal(gRpcServer.Serve(lis))
}