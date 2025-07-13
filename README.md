# url-shortener
## Build Plan

### Day 3 – Resolve & Caching Layer
- Add gRPC unary interceptors for structured logging

### Day 5 – Terraform Infra & CI Integration
- Flesh out Terraform modules:
  - **network**: VPC, subnets, security groups  
  - **database**: RDS Postgres  
  - **cache**: ElastiCache Redis
  - **compute**: ECS cluster, task definitions, IAM roles  
- Add GitHub Actions steps:
  1. `terraform fmt` → `terraform plan` (manual approval) → `terraform apply`

### Day 6 – Deploy & Smoke Test
- Run `terraform apply` to provision prod infra
- Build & push Docker image to ECR
- Deploy ECS service on Fargate with NLB for port 50051

### Day 7 – Gateway, Monitoring & Polish
- Provision CloudWatch alarms via Terraform
- Finalize CI:
  - On merge to `main`: push image, `terraform apply`
- Write README sections:
  - Architecture diagram
  - Usage examples (Shrink, Resolve)

# URL Shortener

---

## 1. Problem Statement & Goals

* **Goal:** Provide planet‑scale URL redirection with <5 ms p95 latency.
* **Non‑Goals:** 

## 2. High‑Level Architecture

*(Insert a real diagram here)*

### 2.1 Components

| Component             | Purpose                               | Tech Stack                       |
| --------------------- | ------------------------------------- | -------------------------------- |

### 2.2 URL Generation & Collision Handling

I opted for a hash + salting system to generate codes for the URLs. Each URL becomes an 8 character Base62 encoded string, giving about 218 trillion unique combinations. Codes have to be unique, so in the event that a collision occurs the code is regenerated with a different salt, ensuring each URL is able to generate a unique code. Users must specify a TTL for their shortened URL of 1 hour, 24 hours, or 168 hours (1 week). 

## 3. API Design

### 3.1 Public HTTP Endpoints

| Method | Path            | Resp            | Notes       |
| ------ | --------------- | --------------- | ----------- |
|  POST  |   `/shorten`    | `{code: string}`| Accepts JSON `{ url: "..."}` |
|  GET   |    `/{code}`    | Redirect (302)  | Looks up code and 302→original URL |

### 3.2 Internal gRPC Services

All of the business logic is exposed over gRPC via the `shortener.Shortener` service. The full definition lives in [`proto/shortener.proto`](proto/shortener.proto):

```proto
message ShortenRequest {
  string url = 1;
}

message ShortenResponse {
  string code = 1;
}

message ResolveRequest {
  string code = 1;
}

message ResolveResponse {
  string url = 1;
}

service Shortener {
  rpc Shorten(ShortenRequest) returns (ShortenResponse);
  rpc Resolve(ResolveRequest) returns (ResolveResponse);
}
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

## 5. Scaling & Capacity Planning

## 6. Consistency & Caching Strategy

## 7. Rate Limiting & Abuse Prevention

## 8. Security Considerations

## 9. Observability

slug_generation_collision_rate

## 10. Deployment & CI/CD

## 11. Testing Strategy

## 12. Failure Modes & Mitigations

---

*Last updated: 2025‑05‑29*
