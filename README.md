# URL Shortener

**🔗 Live Demo: https://frontend-five-lac-36.vercel.app/**

---

## 1. Problem Statement & Goals

* **Goal:** Provide high‑scale URL redirection with <200 ms p95 latency.

## 2. High‑Level Architecture

*(Insert a real diagram here)*

### 2.1 Components

| Component             | Purpose                               | Tech Stack                       |
| --------------------- | ------------------------------------- | -------------------------------- |
| ALB  | TLS termination, HTTP → gRPC routing, health checks    |             AWS ALB              |
|      ECS Service      |   Runs containerised Go application   |          AWS ECS Fargate         |
|  Shortener container  | REST & gRPC API, slug generation, caching layer | Go 1.22, Docker        |
|       PostgreSQL      |       Source‑of‑truth link store      |       Amazon RDS (Postgres)      |
|         Redis         |     Hot‑path cache for code → URL     |        ElastiCache Redis 7       |

### 2.2 URL Generation & Collision Handling

Every long URL is normalised, salted, and hashed. The first 48 bits are then Base‑62 encoded, yielding an 8‑character slug (≈ 218 trillion combinations). If a collision is detected the service re‑salts and retries. Links expire after 24 h by default (TTL stored in expires_at).

## 3. API Design

### 3.1 Public HTTP Endpoints

| Method | Path            | Resp            | Notes       |
| ------ | --------------- | --------------- | ----------- |
|  POST  |   `/shorten`    | `{code: string}`| Accepts JSON `{ url: "..."}` |
|  GET   |    `/{code}`    | Redirect (302)  | Looks up code and 302→original URL |

### 3.2 Internal gRPC Services

All of the business logic is exposed over gRPC via the `shortener.Shortener` service. The full definition lives in [`proto/shortener.proto`](proto/shortener.proto):

```proto
message ShortenRequest { string url = 1; }
message ShortenResponse { string code = 1; }
message ResolveRequest { string code = 1; }
message ResolveResponse { string url = 1; }
service Shortener {
  rpc Shorten(ShortenRequest) returns (ShortenResponse);
  rpc Resolve(ResolveRequest) returns (ResolveResponse);
}
```

## Usage

### Shorten a URL over HTTP

```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"url":"https://example.com/some/very/long/path"}' \
     https://<ALB‑DNS>/shorten
# → {"code":"A7f3eG9b"}
```

### Resolve / follow redirect
```bash
curl -I https://<ALB‑DNS>/A7f3eG9b
# HTTP/1.1 302 Found
# Location: https://example.com/some/very/long/path
```

### Shorten
```bash
grpcurl -plaintext -d '{"url":"https://example.com/some/very/long/path"}' \
  <ALB‑DNS>:50051 shortener.Shortener/Shorten
# { "code": "A7f3eG9b" }
```

### Resolve
```bash
grpcurl -plaintext -d '{"code":"A7f3eG9b"}' \
  <ALB‑DNS>:50051 shortener.Shortener/Resolve
# { "url": "https://example.com/some/very/long/path" }
```

## 4. Data Model

```
CREATE TABLE Links (
  code VARCHAR PK,
  long_url TEXT,
  created_at TIMESTAMP,
  expires_at TIMESTAMP
);
```

## 5. Consistency & Caching Strategy

1. POST /shorten
  a. Begin DB txn → INSERT … ON CONFLICT.
  b. Commit.
  c. Async: SETEX code url TTL in Redis.

2. GET /{code}
  a. GET code from Redis.
  b. Hit → 302 redirect.
  c. Miss → query Postgres, then SETEX with residual TTL.

### Properties

- Read‑after‑write consistency for practically all requests because the writer populates Redis before the first redirect occurs.
- Expiry: Redis TTL = min(link_TTL, 24 h); nightly job purges expired DB rows (DELETE WHERE expires_at ≤ now()).

## 6. Observability

The URL Shortener exposes Prometheus metrics on the `/metrics` endpoint. I scraped these metrics with Prometheus and built dashboards in Grafana to monitor service health and performance.

- shortener_collision_rate
  - Tracks the rate of hash collisions encountered when generating new codes.   
- resolve_cache_hits_total
  - Cumulative count of successful cache lookups during code resolution.
- resolve_cache_misses_total
  - Cumulative count of cache misses when resolving codes.
- resolve_cache_errors_total
  - Total number of errors encountered in the resolve cache layer.

## 7. Deployment & CI/CD

GitHub Actions is used for CI, where linting, unit testing, integration testing, and load testing is done. Afterwards passing the build and tests, the image is uploaded to Docker Hub for use in deployment later. Terraform is then used for provisioning infrastructure and deploying the application to AWS.

I use GitHub Actions to provide a fully automated CI/CD pipeline:

1. **On every push or PR** to `main`/`master`:
   - **Lint** with `golangci-lint`  
   - **Unit tests** (`go test -short`)  
   - **Integration tests** (`go test -tags=integration` against a Docker-Compose stack)  
   - **Load tests** with k6 (smoke-level RPS in CI)  

2. **Build & tag**  
   - `docker compose build app`  
   - Re-tag the image as `mrgoosey/url-shortener:latest`  

3. **Push**  
   - Log in to Docker Hub using secrets  
   - Push the already-tested image  

4. **Terraform-driven deploy**  
   - Plan & apply Terraform modules for VPC, RDS, ElastiCache, ECS, IAM, etc.  
   - Point the ECS task definition at the newly pushed image tag  
   - Perform a rolling update on the service  

5. **Smoke checks**  
   - After deployment, run a brief health-check suite (HTTP status, gRPC ping)  

Secrets like database credentials and Docker Hub tokens are stored in GitHub Actions secrets. Terraform deployment is only done manually, but is triggered as a Github Actions workflow.

## 8. Testing Strategy

I cover three layers of testing:

1. **Unit Tests**  
   - Target all pure-Go logic in `internal/service` (URL validation, Base62 encoding, collision handling)   

2. **Integration Tests**  
   - Spin up a real Postgres + Redis + the Go server in Docker Compose  
   - Drive the gRPC `Shorten` & `Resolve` methods end-to-end  
   - Validate round-trips, error cases, cache-miss vs cache-hit behavior  

3. **Load Tests**  
   - Run k6 scripts against the live HTTP and gRPC endpoints  
   - Measure latency histograms (p50, p95, p99) under realistic RPS  
   - Gate only on broad smoke-level RPS in CI; full performance runs happen in a dedicated environment  
