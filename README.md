# url-shortener
## Build Plan

### Day 3 – Resolve & Caching Layer
- Add gRPC unary interceptors for structured logging and basic rate-limiting

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