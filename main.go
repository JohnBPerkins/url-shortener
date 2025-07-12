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
	dsn := os.Getenv("DATABASE_DSN")
	dbPool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to Postgres (%q): %v", dsn, err)
	}

	//init cache
	cache := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ENDPOINT"),
        Password: "", // if you’re using Redis AUTH
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
	mux.HandleFunc("/api/shorten", shrinkHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        if path == "/" || strings.HasPrefix(path, "/static/") {
            http.FileServer(http.Dir("./web/")).ServeHTTP(w, r)
            return
        }

        code := strings.Trim(path, "/")
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