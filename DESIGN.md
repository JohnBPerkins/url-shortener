# URL Shortener — DESIGN.md Template

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

I opted for a hash + salting system to generate codes for the URLs. Each URL becomes an 6 character Base62 encoded string, giving about 56.8 billion unique combinations. Codes have to be unique, so in the event that a collision occurs the code is regenerated with a different salt, ensuring each URL is able to generate a unique code. Users must specify a TTL for their shortened URL of 1 hour, 24 hours, or 168 hours (1 week). 

## 3. API Design

### 3.1 Public HTTP Endpoints

| Method | Path            | Resp            | Notes       |
| ------ | --------------- | --------------- | ----------- |

### 3.2 Internal gRPC Services

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
