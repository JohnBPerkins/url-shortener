package urlshortener

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/JohnBPerkins/url-shortener/gen"
	"github.com/JohnBPerkins/url-shortener/modules/db"
	"github.com/JohnBPerkins/url-shortener/modules/flake"
	"github.com/JohnBPerkins/url-shortener/service"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func main() {
	//init db
	ctx := context.Background()

	dbPool, _ := db.NewPool(ctx, os.Getenv("DATABASE_URL"))

	//init cache
	cache := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ENDPOINT"), // e.g. "shortener-redis.xxxx.use1.cache.amazonaws.com:6379"
        Password: "", // if you’re using Redis AUTH
        DB:       0,
    })

	flake := flake.NewSonyflake()
	svc := service.NewShortenerService(dbPool, cache, flake)

	gRpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prom.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prom.StreamServerInterceptor),
	)
	
	pb.RegisterShortenerServer(gRpcServer, svc)
	grpc_prom.EnableHandlingTimeHistogram()

	go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Printf("Prometheus metrics listening on :9090/metrics")
        log.Fatal(http.ListenAndServe(":9090", nil))
    }()

	lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Printf("gRPC server listening on %s", lis.Addr())
    log.Fatal(gRpcServer.Serve(lis))
}