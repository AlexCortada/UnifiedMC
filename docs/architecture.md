# Unified IT Operations Portal — Architecture

## Overview

The UOP follows a microservices architecture with a vendor-agnostic connector framework.
Multiple backend systems (Ivanti, Intune, Entra ID) are abstracted behind a unified interface.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  PRESENTATION LAYER (React + TypeScript)                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                     │
│  │ Executive    │  │ Operations   │  │ Remote Ops   │                     │
│  │ Dashboard    │  │ Console      │  │ Center       │                     │
│  └──────────────┘  └──────────────┘  └──────────────┘                     │
├─────────────────────────────────────────────────────────────────────────────┤
│  API GATEWAY (Kong / Envoy)                                                 │
│  • Authentication (JWT)  • Rate Limiting  • Tenant Routing                │
├─────────────────────────────────────────────────────────────────────────────┤
│  CORE SERVICES (Go + TypeScript)                                            │
│  ┌─────────────┐ ┌──────────────┐ ┌─────────────┐ ┌───────────────┐      │
│  │ Inventory   │ │ Merge        │ │ Remote Ops  │ │ Dashboard     │      │
│  │ Service     │ │ Engine       │ │ Service     │ │ Service       │      │
│  └─────────────┘ └──────────────┘ └─────────────┘ └───────────────┘      │
│  ┌─────────────┐ ┌──────────────┐ ┌─────────────┐ ┌───────────────┐      │
│  │ Remediation │ │ Workflow     │ │ Audit       │ │ Notification  │      │
│  │ Engine      │ │ Orchestrator │ │ Service     │ │ Service       │      │
│  └─────────────┘ └──────────────┘ └─────────────┘ └───────────────┘      │
├─────────────────────────────────────────────────────────────────────────────┤
│  CONNECTOR FRAMEWORK                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Connector SDK: Auth, Rate Limit, Circuit Breaker, Retry           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │ Ivanti   │ │ Intune   │ │ Entra ID │ │ ServiceN.│ │ N-Able   │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
├─────────────────────────────────────────────────────────────────────────────┤
│  DATA LAYER                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐        │
│  │ PostgreSQL   │  │ Redis        │  │ Elasticsearch            │        │
│  │ (Primary DB) │  │ (Cache/Sess) │  │ (Search & Analytics)     │        │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐        │
│  │ Kafka        │  │ ClickHouse   │  │ HashiCorp Vault          │        │
│  │ (Events)     │  │ (Time-series)│  │ (Secrets)                │        │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### Frontend
- **Executive Dashboard**: 10 KPI widgets, trend charts, risk heat map
- **Operations Console**: Device list, search, bulk actions
- **Remote Ops Center**: Script editor, target selector, output viewer, shell terminal

### API Gateway
- Authentication via JWT validation
- Rate limiting per user/tenant
- Request routing to backend services
- Tenant context injection

### Core Services
| Service | Responsibility |
|---|---|
| Inventory Service | Device CRUD, search, aggregation |
| Merge Engine | Deduplication, entity resolution, confidence scoring |
| Remote Ops Service | Action routing, execution, fallback, output capture |
| Dashboard Service | KPI computation, trend aggregation, caching |
| Remediation Engine | Policy evaluation, automated remediation |
| Identity & RBAC | Authentication, authorization, permission evaluation |
| Audit Service | Immutable audit trail, search, export |
| Tenant Service | Tenant management, configuration |

### Connector Framework
- **SDK**: Base classes with auth, rate limiting, circuit breaker, retry
- **Registry**: Auto-discovery of connector implementations
- **Canonical Models**: Normalized data structures across all connectors

### Data Layer
| Store | Purpose |
|---|---|
| PostgreSQL | Primary relational store, RLS for tenant isolation |
| Redis | Hot device cache (5-min TTL), sessions, rate limiting |
| Kafka | Event streaming, async processing |
| Elasticsearch | Full-text search, dashboard aggregations |
| ClickHouse | Time-series metrics, trend data |
| HashiCorp Vault | Secrets, credentials, dynamic rotation |

## Security

- TLS 1.3 for all communications
- mTLS between services
- Entra ID authentication with MFA
- Row-Level Security for tenant isolation
- Immutable audit logging
- OWASP Top 10 protection
