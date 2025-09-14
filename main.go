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
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_DSN") // fallback
		if dsn == "" {
			dsn = "postgres://localhost:5432/urlshortener?sslmode=disable"
		}
	}
	dbPool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to Postgres (%q): %v", dsn, err)
	}

	//init cache
	redisURL := os.Getenv("REDIS_URL")
	redisAddr := os.Getenv("REDIS_ENDPOINT")
	if redisURL != "" {
		// Parse Railway Redis URL: redis://default:password@host:port
		redisAddr = strings.TrimPrefix(redisURL, "redis://")
		parts := strings.Split(redisAddr, "@")
		if len(parts) == 2 {
			redisAddr = parts[1]
		}
	} else if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	cache := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: os.Getenv("REDIS_PASSWORD"), // Railway provides this
        DB:       0,
    })

	flake := flake.NewSonyflake()
	svc := service.NewShortenerService(dbPool, cache, flake)

	gRpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prom.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prom.StreamServerInterceptor),
	)

	shrinkHandler := web.NewShrinkHandler(svc)

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        code := strings.Trim(r.URL.Path, "/")
        if code == "" {
            http.Error(w, "Not Found", http.StatusNotFound)
            return
        }

        resp, err := svc.Resolve(r.Context(), &pb.ResolveRequest{Code: code})
        if err != nil {
            http.NotFound(w, r)
            return
        }

        http.Redirect(w, r, resp.GetUrl(), http.StatusFound)
    })

    go func() {
        metricsMux := http.NewServeMux()
        metricsMux.Handle("/metrics", promhttp.Handler())
        log.Println("▶ metrics listening on :9090/metrics")
        log.Fatal(http.ListenAndServe(":9090", metricsMux))
    }()

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