# url-shortener
## 7-Day Build Plan

### Day 1 – Design & Terraform Bootstrap
- ~~Draft `shortener.proto` with services:~~
  - ~~`Shortener.Shrink( url ) returns ( code )`~~
  - ~~`Shortener.Resolve( code ) returns ( url )`~~
- ~~Initialize a Terraform workspace~~
- ~~Create empty module stubs:~~
  - ~~`modules/network/`~~
  - ~~`modules/database/`~~
  - ~~`modules/cache/`~~
  - ~~`modules/compute/`~~

### Day 2 – gRPC Service Core
- ~~Generate Go bindings from `shortener.proto`~~
- ~~Implement `Shrink` RPC:~~
  - ~~URL validation~~
  - ~~Sonyflake/UUID code generation~~
  - ~~Insert into PostgreSQL~~
- Write unit tests for collision handling and invalid URLs

### Day 3 – Resolve & Caching Layer
- Implement `Resolve` RPC:
  - Check Redis cache → fallback to Postgres
- Add gRPC unary interceptors for structured logging and basic rate-limiting
- Stand up a Docker Compose stack: Go service + Postgres + Redis

### Day 4 – Metrics
- Expose Prometheus metrics via a gRPC interceptor

### Day 5 – Terraform Infra & CI Integration
- Flesh out Terraform modules:
  - **network**: VPC, subnets, security groups  
  - **database**: RDS Postgres  
  - **cache**: ElastiCache Redis  
  - **compute**: ECS cluster, task definitions, IAM roles  
- Add GitHub Actions steps:
  1. `terraform fmt` → `terraform plan` (manual approval) → `terraform apply`
  2. `protoc` codegen → `go test` → `golangci-lint`

### Day 6 – Deploy & Smoke Test
- Run `terraform apply` to provision prod infra
- Build & push Docker image to ECR
- Deploy ECS service on Fargate with NLB for port 50051
- Run end-to-end smoke tests (via `grpcurl` or a tiny Go/TS client)

### Day 7 – Gateway, Monitoring & Polish
- (Optional) Add gRPC-Gateway & Envoy for REST façade and TLS
- Provision Grafana dashboards or CloudWatch alarms via Terraform
- Finalize CI:
  - On merge to `main`: lint, test, build, push image, `terraform apply`
- Write README sections:
  - Architecture diagram
  - Proto reference
  - Usage examples (Shrink, Resolve, Stats)
