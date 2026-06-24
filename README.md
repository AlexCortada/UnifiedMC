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

## Development Workflow

```
1. Create/update files locally
2. Push to GitHub repo (source of truth)
3. SSH to Mission Control (UMC): ssh hermes@10.0.10.123
4. Pull updates: cd /opt/unifiedmc && git pull origin main
5. Validate: run tests, check services, verify API
```

### Mission Control (UMC)
- **Host:** 10.0.10.123
- **OS:** Ubuntu 24.04 LTS
- **User:** hermes
- **App Directory:** /opt/unifiedmc
- **Specs:** 16GB RAM, 256GB SSD, i5

### GitHub Repo
- **URL:** https://github.com/AlexCortada/UnifiedMC/
- **Branch:** main (all changes merged here)
- **Purpose:** Source of truth, backup, collaboration

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
