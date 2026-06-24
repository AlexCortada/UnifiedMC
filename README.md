# Unified IT Operations Portal (UOP)

A vendor-agnostic platform that consolidates device management, application deployment, patch management, remediation, and remote operations across Ivanti Neurons, Microsoft Intune, Microsoft Entra ID, and more.

## Project Specification

The full Project Specification Document (PSD) is maintained at:
- `/docs/Project-Specification.md` (source of truth for all development)

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  Frontend (React + TypeScript)                                  │
│  • Executive Dashboard  • Operations Console  • Remote Shell    │
├─────────────────────────────────────────────────────────────────┤
│  API Gateway (Kong / Envoy)                                     │
│  • Auth (JWT)  • Rate Limiting  • Tenant Routing                │
├─────────────────────────────────────────────────────────────────┤
│  Core Services (Go + TypeScript)                                │
│  • Inventory  • Merge Engine  • Ops  • Dashboard  • Audit       │
├─────────────────────────────────────────────────────────────────┤
│  Connector Framework                                            │
│  • Ivanti Neurons  • Microsoft Intune  • Entra ID  • ServiceNow │
├─────────────────────────────────────────────────────────────────┤
│  Data Layer                                                     │
│  • PostgreSQL  • Redis  • Kafka  • Elasticsearch  • ClickHouse  │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for backend services)

### Local Development

```bash
# Start data infrastructure
docker compose -f infrastructure/docker/docker-compose.yml up -d

# Run database migrations
for f in database/migrations/*.sql; do
  PGPASSWORD=*** psql -h localhost -U unifiedmc -d unifiedmc -f "$f"
done

# Start API gateway
cd services/api-gateway && go run .
```

## Project Structure

```
.
├── services/          # Backend microservices (Go)
│   ├── api-gateway/   # REST/GraphQL gateway
│   ├── inventory-service/
│   ├── merge-engine/
│   ├── dashboard-service/
│   ├── ops-service/   # Remote operations
│   ├── identity-service/
│   └── audit-service/
├── connectors/        # Vendor connector implementations
│   ├── sdk/           # Connector SDK (base classes)
│   ├── intune/        # Microsoft Intune connector
│   ├── ivanti/        # Ivanti Neurons connector
│   ├── servicenow/    # ServiceNow connector
│   └── nable/         # N-Able connector
├── web/               # React frontend
├── packages/          # Shared libraries
│   ├── types/         # Canonical data models
│   └── config/        # Shared configuration
├── infrastructure/    # Deployment configs
│   ├── docker/
│   ├── k8s/
│   └── terraform/
├── database/          # Migrations & seeds
│   ├── migrations/
│   └── seeds/
├── docs/              # Architecture & specs
└── .github/workflows/ # CI/CD
    └── ci.yml
```

## Phases

| Phase | Timeline | Deliverables |
|-------|----------|--------------|
| Phase 0 | Days 1-3 | Scaffolding, CI/CD, Docker Compose |
| Phase 1 | Days 4-7 | Database schema, Redis, seed data |
| Phase 2 | Days 8-12 | Connector SDK, base class, mock connector |
| Phase 3 | Days 13-18 | Intune connector (read-only) |
| Phase 4 | Days 19-23 | Unified Device API |
| Phase 5 | Days 24-30 | Ivanti connector + Merge engine |
| Phase 6 | Days 31-37 | Dashboard UI + Device list + Auth |

## License

See [LICENSE](./LICENSE).
