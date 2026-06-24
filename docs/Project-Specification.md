# Project Specification Document (PSD)

## Unified IT Operations Portal

**Document Version:** 1.0  
**Classification:** Internal — Source of Truth  
**Date:** 2026-06-18  
**Status:** Baseline  
**Author:** Enterprise Architecture  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)  
2. [Functional Requirements](#2-functional-requirements)  
3. [Non-Functional Requirements](#3-non-functional-requirements)  
4. [User Roles](#4-user-roles)  
5. [System Architecture](#5-system-architecture)  
6. [Database Schema](#6-database-schema)  
7. [API Design](#7-api-design)  
8. [Connector Framework](#8-connector-framework)  
9. [Security Model](#9-security-model)  
10. [Deployment Model](#10-deployment-model)  
11. [MVP Scope](#11-mvp-scope)  
12. [Future Scope](#12-future-scope)  
13. [Appendices](#13-appendices)

---

## 1. Executive Summary

### 1.1 Purpose

The Unified IT Operations Portal (UOP) consolidates device management, application deployment, patch management, remediation, and remote operations across multiple vendor platforms into a single pane of glass. It eliminates tool fragmentation by providing vendor-agnostic connectors, a canonical data model, and unified workflows that route actions to the best available backend without requiring end users to know which platform executes their requests.

### 1.2 Scope

This document defines the complete architecture, requirements, and implementation plan for the UOP. It is the authoritative source of truth for all development teams.

### 1.3 Key Objectives

| Objective | Description |
|---|---|
| Unified Inventory | Single canonical device view across Ivanti, Intune, Entra ID |
| Unified Deployment | Application deployment routed to best available backend |
| Unified Remediation | Policy-based automated remediation across all platforms |
| Unified Patch Management | Aggregated patch compliance with cross-platform deployment |
| Unified Remote Actions | Script execution, restart, shell — backend-agnostic |
| Executive Dashboards | Real-time KPIs covering devices, compliance, vulnerabilities, SLA, HR |

### 1.4 Integrations

| System | Type | Status |
|---|---|---|
| Ivanti Neurons | ITSM / ITAM / Patch Management | Phase 1 |
| Microsoft Intune | MDM / MAM | Phase 1 |
| Microsoft Entra ID | Identity Provider | Phase 1 |
| ServiceNow | ITSM / CMDB | Phase 2 |
| Microsoft Defender | Endpoint Security | Phase 4 |
| N-Able | RMM | Phase 4 |
| Addigy | Apple MDM | Phase 4 |
| Power BI | Analytics / Reporting | Phase 4 |

---

## 2. Functional Requirements

### 2.1 Unified Device Inventory

| ID | Requirement | Priority |
|---|---|---|
| FR-INV-001 | The system shall ingest device records from Ivanti Neurons, Microsoft Intune, and Microsoft Entra ID. | P0 |
| FR-INV-002 | The system shall merge records from multiple sources into a single canonical device object using configurable matching rules. | P0 |
| FR-INV-003 | Matching rules shall cascade: serial number (highest confidence), MAC address, hostname + OS, user + device name, IP + time window. | P0 |
| FR-INV-004 | The system shall compute a confidence score (0.0–1.0) for each merge. | P0 |
| FR-INV-005 | Records with confidence ≥ 0.90 shall be auto-merged. Records with 0.70–0.89 shall be flagged for manual review. | P0 |
| FR-INV-006 | The system shall display the following fields per device: Device Name, User, Department, OS, Compliance Status, Patch Status, Last Seen, Installed Applications. | P0 |
| FR-INV-007 | The system shall support manual merge and split operations for duplicate management. | P1 |
| FR-INV-008 | The system shall maintain full provenance showing which backend sources contributed to each device record. | P0 |
| FR-INV-009 | The system shall support full-text search across device name, user, department, and installed applications. | P1 |
| FR-INV-010 | The system shall track installed software per device with version, publisher, install date, and approval status. | P1 |
| FR-INV-011 | The system shall support device grouping by OS type, compliance status, department, organization, and custom tags. | P1 |
| FR-INV-012 | The system shall compute a risk score (0–100) per device based on compliance, patch status, vulnerabilities, and activity. | P2 |

### 2.2 Application Deployment

| ID | Requirement | Priority |
|---|---|---|
| FR-DEP-001 | The system shall deploy applications to devices through Ivanti, Intune, or N-Able based on automatic backend selection. | P0 |
| FR-DEP-002 | The system shall support single-device and bulk deployment (up to 10,000 devices). | P0 |
| FR-DEP-003 | The system shall track deployment status per device: pending, running, completed, failed, cancelled. | P0 |
| FR-DEP-004 | The system shall support scheduled deployments with date/time windows. | P1 |
| FR-DEP-005 | The system shall support deployment rollback for failed deployments. | P2 |
| FR-DEP-006 | The system shall maintain an application catalog with approval workflow. | P1 |

### 2.3 Patch Management

| ID | Requirement | Priority |
|---|---|---|
| FR-PAT-001 | The system shall aggregate patch status from Ivanti and Intune into a unified compliance view. | P0 |
| FR-PAT-002 | The system shall display patch compliance by severity (critical, high, medium, low). | P0 |
| FR-PAT-003 | The system shall support deploying patches to devices through the best available backend. | P0 |
| FR-PAT-004 | The system shall support patch policies that define approval rules and deployment windows. | P1 |
| FR-PAT-005 | The system shall track patch installation history per device. | P1 |
| FR-PAT-006 | The system shall identify superseded patches and exclude them from compliance calculations. | P2 |

### 2.4 Remote Operations

| ID | Requirement | Priority |
|---|---|---|
| FR-OPS-001 | The system shall execute scripts (PowerShell, Bash, Python) on target devices through Ivanti, Intune, or N-Able. | P0 |
| FR-OPS-002 | The system shall restart/reboot devices through the best available backend. | P0 |
| FR-OPS-003 | The system shall provide interactive remote shell sessions through Ivanti or N-Able. | P1 |
| FR-OPS-004 | The system shall automatically select the best backend for each action based on device capabilities, source health, freshness, and admin preferences. | P0 |
| FR-OPS-005 | The system shall support automatic fallback to secondary backends if the primary fails. | P0 |
| FR-OPS-006 | The system shall support bulk actions on up to 1,000 devices simultaneously. | P1 |
| FR-OPS-007 | The system shall stream real-time script output via WebSocket. | P1 |
| FR-OPS-008 | The system shall record remote shell sessions for compliance auditing. | P2 |
| FR-OPS-009 | The system shall maintain a script library with parameterization and approval workflow. | P1 |
| FR-OPS-010 | The system shall support action scheduling with approval workflows. | P2 |

### 2.5 Remediation Engine

| ID | Requirement | Priority |
|---|---|---|
| FR-REM-001 | The system shall detect compliance violations across all connected platforms. | P1 |
| FR-REM-002 | The system shall support automated remediation policies with configurable triggers and actions. | P1 |
| FR-REM-003 | The system shall support manual approval for remediation actions based on risk level. | P1 |
| FR-REM-004 | The system shall track remediation history and success rates. | P2 |
| FR-REM-005 | The system shall support rollback of remediation actions. | P2 |

### 2.6 Executive Dashboard

| ID | Requirement | Priority |
|---|---|---|
| FR-DASH-001 | The system shall display the following KPIs on a single dashboard: Total Devices, Online Devices, Offline Devices, Compliance Rate, Patch Compliance, Critical Vulnerabilities, Open Incidents, SLA Breaches, New Hires, Terminations. | P0 |
| FR-DASH-002 | The system shall support date range filtering (24h, 7d, 30d, 90d, custom). | P0 |
| FR-DASH-003 | The system shall display trend charts for compliance, patch status, and device count over time. | P1 |
| FR-DASH-004 | The system shall display device distribution by OS type and department. | P1 |
| FR-DASH-005 | The system shall display a risk heat map by organizational unit. | P2 |
| FR-DASH-006 | The system shall display a real-time activity feed of recent events. | P1 |
| FR-DASH-007 | The system shall support CSV and PDF export of dashboard data. | P2 |
| FR-DASH-008 | The system shall auto-refresh dashboard data every 60 seconds. | P1 |
| FR-DASH-009 | The system shall compute an overall health status (healthy/warning/critical) from KPI thresholds. | P1 |

### 2.7 Audit & Reporting

| ID | Requirement | Priority |
|---|---|---|
| FR-AUD-001 | The system shall log all actions (create, update, delete, execute) in an immutable audit trail. | P0 |
| FR-AUD-002 | The system shall record which backend executed each action. | P0 |
| FR-AUD-003 | The system shall support audit log search by actor, action, resource, date range. | P1 |
| FR-AUD-004 | The system shall retain audit logs for 7 years (compliance requirement). | P0 |
| FR-AUD-005 | The system shall support custom report generation. | P2 |

---

## 3. Non-Functional Requirements

### 3.1 Performance

| ID | Requirement | Target |
|---|---|---|
| NFR-PERF-001 | Device list query (page of 50) response time | < 100ms (p95) |
| NFR-PERF-002 | Device detail query response time | < 50ms (p95) |
| NFR-PERF-003 | Full-text search response time | < 200ms (p95) |
| NFR-PERF-004 | Dashboard summary API response time | < 200ms (p95) |
| NFR-PERF-005 | Action execution per device (excluding backend latency) | < 500ms |
| NFR-PERF-006 | Sync processing per device (map + merge) | < 200ms |
| NFR-PERF-007 | Concurrent device capacity per tenant | 100,000 devices |
| NFR-PERF-008 | Concurrent tenants | 50 |
| NFR-PERF-009 | Bulk action capacity | 1,000 devices per job |
| NFR-PERF-010 | Dashboard auto-refresh interval | 60 seconds |

### 3.2 Availability & Reliability

| ID | Requirement | Target |
|---|---|---|
| NFR-AVAIL-001 | System uptime SLA | 99.9% (8.76 hours downtime/year) |
| NFR-AVAIL-002 | Recovery Point Objective (RPO) | 0 (synchronous replication for critical data) |
| NFR-AVAIL-003 | Recovery Time Objective (RTO) | < 30 minutes |
| NFR-AVAIL-004 | Single backend failure shall not prevent actions from executing on other backends | Required |
| NFR-AVAIL-005 | Circuit breaker shall prevent cascade failures when a backend is unreachable | Required |

### 3.3 Scalability

| ID | Requirement | Target |
|---|---|---|
| NFR-SCAL-001 | Horizontal scaling for API tier | Auto-scale at 70% CPU |
| NFR-SCAL-002 | Horizontal scaling for sync workers | Auto-scale by queue depth |
| NFR-SCAL-003 | Database connection pooling | PgBouncer, max 200 connections |
| NFR-SCAL-004 | Read replicas for query distribution | Minimum 2 read replicas |
| NFR-SCAL-005 | Elasticsearch scaling | 3 primary shards, 1 replica |

### 3.4 Security

| ID | Requirement | Target |
|---|---|---|
| NFR-SEC-001 | All API communications shall use TLS 1.3 | Required |
| NFR-SEC-002 | Service-to-service communication shall use mTLS | Required |
| NFR-SEC-003 | Authentication via Microsoft Entra ID (SAML 2.0 / OIDC) | Required |
| NFR-SEC-004 | Multi-factor authentication enforced at identity provider | Required |
| NFR-SEC-005 | RBAC with tenant isolation at every layer | Required |
| NFR-SEC-006 | Data at rest encrypted (AES-256) | Required |
| NFR-SEC-007 | Secrets stored in HashiCorp Vault | Required |
| NFR-SEC-008 | PII redaction in logs and API responses | Required |
| NFR-SEC-009 | Rate limiting per API key / user | Required |
| NFR-SEC-010 | OWASP Top 10 protection | Required |

### 3.5 Maintainability

| ID | Requirement | Target |
|---|---|---|
| NFR-MAINT-001 | New connector integration shall require only implementing the interface + registration | Required |
| NFR-MAINT-002 | Zero-downtime deployments (blue/green or canary) | Required |
| NFR-MAINT-003 | All services shall expose OpenTelemetry metrics | Required |
| NFR-MAINT-004 | Structured logging (JSON) with correlation IDs | Required |
| NFR-MAINT-005 | Database migrations shall be reversible | Required |

### 3.6 Compliance

| ID | Requirement | Target |
|---|---|---|
| NFR-COMP-001 | SOC 2 Type II readiness | Required |
| NFR-COMP-002 | GDPR data handling (right to erasure, data portability) | Required |
| NFR-COMP-003 | Audit log immutability | Required |
| NFR-COMP-004 | Data retention: audit logs 7 years, device data life of tenant | Required |

---

## 4. User Roles

### 4.1 Role Definitions

| Role | Description | Access Scope |
|---|---|---|
| **Platform Admin** | Full system access, tenant management, connector configuration | All tenants, all data |
| **Tenant Admin** | Tenant-level administration, user management, preferences | Single tenant, all data |
| **Security Admin** | Security policy management, vulnerability management, audit review | Single tenant, security data |
| **Operations Manager** | Approve high-risk actions, view all operations | Single tenant, operational data |
| **Deployment Engineer** | Create and execute deployments, manage application catalog | Single tenant, deployment data |
| **Patch Engineer** | Create patch policies, execute patch deployments | Single tenant, patch data |
| **Service Desk Analyst** | Execute low-risk actions (restart, scripts), view devices | Single tenant, assigned org |
| **Compliance Auditor** | Read-only access to compliance data, audit logs | Single tenant, read-only |
| **Read-Only Viewer** | View dashboards and device data, no actions | Single tenant, read-only |
| **End User** | Self-service portal: view own devices, request actions | Own devices only |

### 4.2 Permission Matrix

| Capability | Platform Admin | Tenant Admin | Security Admin | Ops Mgr | Deploy Eng | Patch Eng | Service Desk | Auditor | Viewer | End User |
|---|---|---|---|---|---|---|---|---|---|---|
| View Dashboard | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 🔶 |
| View Devices | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 🔶 |
| Search Devices | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 🔶 |
| Run Script | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | 🔶 | ❌ | ❌ | ❌ |
| Restart Device | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ |
| Deploy Application | ✅ | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Deploy Patch | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Remote Shell | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | 🔶 | ❌ | ❌ | ❌ |
| Bulk Actions (>50) | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Approve Actions | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| View Audit Log | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| Manage Connectors | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Manage Users/Roles | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Configure Policies | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Export Data | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 🔶 | ✅ | 🔶 | ❌ |
| Request Action | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

**Legend:** ✅ Full access | 🔶 Restricted scope | ❌ No access

---

## 5. System Architecture

### 5.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    UNIFIED IT OPERATIONS PORTAL                              │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  PRESENTATION LAYER                                                   │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │  │
│  │  │ Executive    │  │ Operations   │  │ Remote Ops   │               │  │
│  │  │ Dashboard    │  │ Console      │  │ Center       │               │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│  ┌───────────────────────────▼───────────────────────────────────────────┐  │
│  │  API GATEWAY (Kong / Envoy)                                           │  │
│  │  • Authentication (JWT validation)                                    │  │
│  │  • Rate limiting (per user / tenant)                                  │  │
│  │  • Request validation                                                 │  │
│  │  • Tenant context injection                                          │  │
│  └───────────────────────────┬───────────────────────────────────────────┘  │
│                              │                                              │
│  ┌───────────────────────────▼───────────────────────────────────────────┐  │
│  │  CORE SERVICES (Kubernetes)                                           │  │
│  │                                                                       │  │
│  │  ┌─────────────┐ ┌──────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Inventory   │ │ Merge        │ │ Remote Ops  │ │ Dashboard     │  │  │
│  │  │ Service     │ │ Engine       │ │ Service     │ │ Service       │  │  │
│  │  └─────────────┘ └──────────────┘ └─────────────┘ └───────────────┘  │  │
│  │  ┌─────────────┐ ┌──────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Remediation │ │ Workflow     │ │ Audit       │ │ Notification  │  │  │
│  │  │ Engine      │ │ Orchestrator │ │ Service     │ │ Service       │  │  │
│  │  └─────────────┘ └──────────────┘ └─────────────┘ └───────────────┘  │  │
│  │  ┌─────────────┐ ┌──────────────┐                                     │  │
│  │  │ Identity &  │ │ Tenant       │                                     │  │
│  │  │ RBAC Svc    │ │ Service      │                                     │  │
│  │  └─────────────┘ └──────────────┘                                     │  │
│  └───────────────────────────┬───────────────────────────────────────────┘  │
│                              │                                              │
│  ┌───────────────────────────▼───────────────────────────────────────────┐  │
│  │  CONNECTOR FRAMEWORK                                                  │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Connector SDK: Auth, Rate Limit, Circuit Breaker, Retry       │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐               │  │
│  │  │ Ivanti   │ │ Intune   │ │ Entra ID │ │ ServiceN.│  (Future: ...)  │  │
│  │  │Connector │ │Connector │ │Connector │ │Connector │               │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│  ┌───────────────────────────▼───────────────────────────────────────────┐  │
│  │  DATA LAYER                                                           │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │  │
│  │  │ PostgreSQL   │  │ Redis        │  │ Elasticsearch            │   │  │
│  │  │ (Primary DB) │  │ (Cache/Sess) │  │ (Search & Analytics)     │   │  │
│  │  └──────────────┘  └──────────────┘  └──────────────────────────┘   │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │  │
│  │  │ Kafka        │  │ ClickHouse   │  │ HashiCorp Vault          │   │  │
│  │  │ (Events)     │  │ (Time-series)│  │ (Secrets)                │   │  │
│  │  └──────────────┘  └──────────────┘  └──────────────────────────┘   │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Technology Stack

| Layer | Technology | Rationale |
|---|---|---|
| Frontend | React 18 + TypeScript, TanStack Query, Tailwind CSS | Component-rich, type-safe, excellent DX |
| API Gateway | Kong / Envoy | Rate limiting, auth, GraphQL federation |
| Core Services | Go (performance) + Python (ML/analytics) | Go for connectors; Python for analytics |
| Workflow Engine | Temporal.io | Durable execution, saga patterns, retry |
| Event Bus | Apache Kafka / NATS JetStream | Durable, ordered, replayable |
| Primary DB | PostgreSQL 16 | Relational integrity, JSONB, RLS |
| Cache | Redis Cluster (6 nodes) | Hot data, sessions, rate limiting |
| Search | Elasticsearch / OpenSearch | Full-text, aggregations |
| Analytics | ClickHouse | Columnar, time-series, fast aggregates |
| Secrets | HashiCorp Vault | Centralized, dynamic credentials |
| Observability | OpenTelemetry + Grafana Stack | Vendor-neutral telemetry |
| Orchestration | Kubernetes (EKS/AKS/GKE) | Auto-scaling, self-healing |
| CI/CD | GitHub Actions + ArgoCD | GitOps, blue/green deploys |

---

## 6. Database Schema

### 6.1 Core Tables

#### `tenants`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| name | VARCHAR(255) | NOT NULL |
| domain | VARCHAR(255) | UNIQUE |
| status | ENUM | `active`, `suspended`, `provisioning` |
| settings | JSONB | |
| created_at | TIMESTAMPTZ | NOT NULL |

#### `organizations`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | FK → tenants, NOT NULL |
| name | VARCHAR(255) | NOT NULL |
| parent_org_id | UUID | FK → organizations |
| path | LTREE | Materialized path for hierarchy |
| status | ENUM | `active`, `inactive` |
| created_at | TIMESTAMPTZ | NOT NULL |

#### `unified_devices`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | FK → tenants, NOT NULL |
| org_id | UUID | FK → organizations |
| display_name | VARCHAR(255) | NOT NULL |
| asset_type | ENUM | `workstation`, `server`, `mobile`, `iot`, `vm`, `container` |
| os_type | ENUM | `windows`, `macos`, `linux`, `ios`, `android`, `other` |
| os_version | VARCHAR(100) | |
| primary_user_id | UUID | FK → device_source_users |
| department | VARCHAR(255) | |
| compliance_status | ENUM | `compliant`, `non_compliant`, `unknown`, `exempt` |
| patch_status_summary | JSONB | `{ "critical_missing": 0, "high_missing": 1 }` |
| last_seen | TIMESTAMPTZ | |
| risk_score | SMALLINT | 0–100 |
| installed_apps_count | INT | |
| merged_from_sources | VARCHAR(50)[] | |
| merge_confidence | FLOAT | 0.0–1.0 |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |
| updated_by_sync_at | TIMESTAMPTZ | |

#### `device_sources`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| unified_device_id | UUID | FK → unified_devices |
| connector_type | VARCHAR(50) | NOT NULL |
| external_id | VARCHAR(500) | NOT NULL |
| raw_data | JSONB | |
| mapped_data | JSONB | |
| last_sync_at | TIMESTAMPTZ | |
| sync_status | ENUM | `synced`, `pending`, `error`, `stale` |
| trust_score | SMALLINT | 0–100 |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

#### `device_source_mappings`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| unified_device_id | UUID | FK → unified_devices |
| connector_type | VARCHAR(50) | |
| external_id | VARCHAR(500) | |
| match_confidence | FLOAT | |
| match_method | VARCHAR(100) | |
| matched_at | TIMESTAMPTZ | |
| matched_by | VARCHAR(100) | |

#### `device_source_users`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | FK → tenants |
| connector_type | VARCHAR(50) | |
| external_id | VARCHAR(500) | |
| email | VARCHAR(255) | |
| display_name | VARCHAR(255) | |
| first_name | VARCHAR(100) | |
| last_name | VARCHAR(100) | |
| department | VARCHAR(255) | |
| job_title | VARCHAR(255) | |
| manager_external_id | VARCHAR(500) | |
| status | VARCHAR(50) | |
| raw_data | JSONB | |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

#### `device_applications`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| unified_device_id | UUID | FK → unified_devices |
| application_name | VARCHAR(255) | |
| application_version | VARCHAR(100) | |
| publisher | VARCHAR(255) | |
| install_date | DATE | |
| install_path | VARCHAR(500) | |
| source_connector | VARCHAR(50) | |
| is_approved | BOOLEAN | |
| raw_data | JSONB | |
| created_at | TIMESTAMPTZ | |

#### `device_patches`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| unified_device_id | UUID | FK → unified_devices |
| patch_name | VARCHAR(255) | |
| kb_article | VARCHAR(50) | |
| severity | ENUM | `critical`, `high`, `medium`, `low` |
| status | ENUM | `installed`, `missing`, `pending`, `failed`, `superseded` |
| source_connector | VARCHAR(50) | |
| detected_at | TIMESTAMPTZ | |
| installed_at | TIMESTAMPTZ | |
| superseded | BOOLEAN | |

### 6.2 Operations Tables

#### `action_requests`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | FK → tenants, NOT NULL |
| action_type | ENUM | `run_script`, `restart_device`, `deploy_application`, `patch_device`, `remote_shell` |
| status | ENUM | `pending`, `validating`, `awaiting_approval`, `approved`, `queued`, `executing`, `completed`, `failed`, `cancelled`, `partial` |
| priority | ENUM | `low`, `normal`, `high`, `critical` |
| target_type | ENUM | `device`, `device_group`, `org` |
| target_ids | UUID[] | NOT NULL |
| target_filter | JSONB | |
| parameters | JSONB | |
| backend_preference | VARCHAR(50) | |
| selected_backend | VARCHAR(50) | |
| fallback_backend | VARCHAR(50) | |
| initiated_by | UUID | FK → users, NOT NULL |
| approved_by | UUID | FK → users |
| approved_at | TIMESTAMPTZ | |
| requires_approval | BOOLEAN | |
| workflow_id | UUID | |
| started_at | TIMESTAMPTZ | |
| completed_at | TIMESTAMPTZ | |
| duration_seconds | INT | |
| result_summary | JSONB | |
| error_details | JSONB | |
| cancellation_reason | TEXT | |
| created_at | TIMESTAMPTZ | NOT NULL |
| updated_at | TIMESTAMPTZ | NOT NULL |

#### `action_targets`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| action_request_id | UUID | FK → action_requests |
| unified_device_id | UUID | FK → unified_devices |
| backend_type | VARCHAR(50) | |
| external_id | VARCHAR(500) | |
| status | ENUM | `pending`, `executing`, `completed`, `failed`, `skipped` |
| result_code | INT | |
| output | TEXT | |
| error | TEXT | |
| started_at | TIMESTAMPTZ | |
| completed_at | TIMESTAMPTZ | |
| duration_seconds | INT | |
| retry_count | INT | |
| raw_result | JSONB | |

#### `action_scripts`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | FK → tenants |
| name | VARCHAR(255) | |
| description | TEXT | |
| script_type | ENUM | `powershell`, `bash`, `python`, `cmd`, `shell` |
| script_content | TEXT | |
| parameters_schema | JSONB | |
| compatible_backends | VARCHAR(50)[] | |
| category | VARCHAR(100) | |
| created_by | UUID | FK → users |
| approved | BOOLEAN | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

#### `remote_shell_sessions`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| action_request_id | UUID | FK → action_requests |
| unified_device_id | UUID | FK → unified_devices |
| backend_type | VARCHAR(50) | |
| session_id | VARCHAR(500) | |
| status | ENUM | `active`, `closed`, `expired`, `terminated` |
| initiated_by | UUID | FK → users |
| started_at | TIMESTAMPTZ | |
| ended_at | TIMESTAMPTZ | |
| duration_seconds | INT | |
| commands_executed | INT | |
| recording_url | VARCHAR(500) | |
| terminated_by | UUID | FK → users |
| termination_reason | TEXT | |

### 6.3 Audit & Reporting Tables

#### `merge_audit_log`

| Column | Type | Constraints |
|---|---|---|
| id | BIGSERIAL | PRIMARY KEY |
| tenant_id | UUID | |
| unified_device_id | UUID | |
| action | VARCHAR(50) | `auto_merge`, `manual_merge`, `split`, `update` |
| source_devices | UUID[] | |
| merge_strategy | VARCHAR(100) | |
| confidence_before | FLOAT | |
| confidence_after | FLOAT | |
| merged_fields | JSONB | |
| performed_by | VARCHAR(100) | |
| performed_at | TIMESTAMPTZ | |

#### `action_audit_log`

| Column | Type | Constraints |
|---|---|---|
| id | BIGSERIAL | PRIMARY KEY |
| tenant_id | UUID | |
| action_request_id | UUID | |
| event_type | VARCHAR(100) | |
| event_data | JSONB | |
| actor_id | UUID | |
| actor_type | VARCHAR(50) | |
| ip_address | INET | |
| timestamp | TIMESTAMPTZ | |

#### `incidents`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | |
| incident_number | VARCHAR(50) | |
| external_id | VARCHAR(500) | |
| source_connector | VARCHAR(50) | |
| title | VARCHAR(500) | |
| description | TEXT | |
| priority | ENUM | `p1`, `p2`, `p3`, `p4` |
| status | ENUM | `new`, `in_progress`, `resolved`, `closed`, `cancelled` |
| category | VARCHAR(100) | |
| assignee_id | UUID | FK → device_source_users |
| department | VARCHAR(255) | |
| created_at | TIMESTAMPTZ | |
| resolved_at | TIMESTAMPTZ | |
| closed_at | TIMESTAMPTZ | |
| sla_response_deadline | TIMESTAMPTZ | |
| sla_resolution_deadline | TIMESTAMPTZ | |
| sla_response_breached | BOOLEAN | |
| sla_resolution_breached | BOOLEAN | |
| response_time_minutes | INT | |
| resolution_time_minutes | INT | |
| raw_data | JSONB | |

#### `vulnerability_findings`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | |
| cve_id | VARCHAR(20) | |
| cvss_score | DECIMAL(3,1) | |
| severity | ENUM | `critical`, `high`, `medium`, `low` |
| product | VARCHAR(255) | |
| vendor | VARCHAR(255) | |
| description | TEXT | |
| patch_available | BOOLEAN | |
| exploit_available | BOOLEAN | |
| published_date | DATE | |
| discovered_at | TIMESTAMPTZ | |
| status | ENUM | `open`, `mitigated`, `accepted`, `false_positive` |
| affected_devices_count | INT | |
| raw_data | JSONB | |

#### `hr_events`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| tenant_id | UUID | |
| event_type | ENUM | `hire`, `termination`, `transfer`, `leave`, `return` |
| user_id | UUID | FK → device_source_users |
| department | VARCHAR(255) | |
| job_title | VARCHAR(255) | |
| event_date | DATE | |
| effective_date | DATE | |
| devices_assigned | INT | |
| devices_recovered | INT | |
| access_revoked | BOOLEAN | |
| onboarding_status | ENUM | `pending`, `in_progress`, `completed` |
| offboarding_status | ENUM | `pending`, `in_progress`, `completed` |
| source_connector | VARCHAR(50) | |
| raw_data | JSONB | |
| created_at | TIMESTAMPTZ | |

### 6.4 Materialized Views

| View | Purpose | Refresh Interval |
|---|---|---|
| `mv_device_availability` | Online/offline counts by tenant | 5 minutes |
| `mv_compliance_summary` | Compliance rates by tenant | 5 minutes |
| `mv_patch_compliance` | Patch compliance by severity | 15 minutes |
| `mv_incident_metrics` | Open incidents, MTTR | 5 minutes |
| `mv_sla_breaches` | Active SLA breaches | 5 minutes |
| `mv_hr_metrics` | Hiring/termination metrics | 1 hour |

---

## 7. API Design

### 7.1 API Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    API ARCHITECTURE                                          │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  GRAPHQL FEDERATION (Primary)                                         │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │  │
│  │  │ Inventory│ │ Ops      │ │ Dashboard│ │ Identity │ │ Audit    │  │  │
│  │  │ Subgraph │ │ Subgraph │ │ Subgraph │ │ Subgraph │ │ Subgraph │  │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  REST API (v1) — For webhooks, exports, file operations              │  │
│  │  /api/v1/devices/*                                                    │  │
│  │  /api/v1/ops/*                                                        │  │
│  │  /api/v1/dashboard/*                                                  │  │
│  │  /api/v1/audit/*                                                      │  │
│  │  /api/v1/connectors/*                                                 │  │
│  │  /api/v1/sync/*                                                       │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  WEBSOCKET API — Real-time updates                                    │  │
│  │  /ws/v1/ops/actions/:id/stream                                        │  │
│  │  /ws/v1/ops/shell/:sessionId                                          │  │
│  │  /ws/v1/dashboard/events                                              │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  WEBHOOK RECEIVER — Inbound from connectors                          │  │
│  │  /webhooks/v1/:connector                                              │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Device Inventory APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/devices` | List devices (paginated, filterable) |
| GET | `/api/v1/devices/:id` | Get unified device details |
| GET | `/api/v1/devices/:id/sources` | Get source records for device |
| GET | `/api/v1/devices/:id/capabilities` | Get device capabilities |
| GET | `/api/v1/devices/:id/history` | Get action history for device |
| POST | `/api/v1/devices/search` | Advanced search |
| POST | `/api/v1/devices/:id/reconcile` | Trigger manual reconciliation |
| POST | `/api/v1/devices/:id/split` | Split merged device |
| POST | `/api/v1/devices/merge` | Manually merge devices |

### 7.3 Operations APIs

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/ops/actions` | Submit action request |
| POST | `/api/v1/ops/actions/bulk` | Submit bulk action |
| GET | `/api/v1/ops/actions/:id` | Get action status |
| GET | `/api/v1/ops/actions/:id/output` | Get action output |
| POST | `/api/v1/ops/actions/:id/cancel` | Cancel action |
| GET | `/api/v1/ops/actions/history` | Action history |
| GET | `/api/v1/ops/scripts` | List script library |
| POST | `/api/v1/ops/scripts` | Add script to library |
| POST | `/api/v1/ops/devices/:id/shell` | Initiate remote shell |
| WS | `/ws/v1/ops/shell/:sessionId` | Shell WebSocket |
| DELETE | `/api/v1/ops/shell/:sessionId` | Terminate shell |
| GET | `/api/v1/ops/preferences` | Get backend preferences |
| PUT | `/api/v1/ops/preferences` | Update backend preferences |

### 7.4 Dashboard APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/dashboard/summary` | All KPIs in one response |
| GET | `/api/v1/dashboard/trends/compliance` | Compliance trend |
| GET | `/api/v1/dashboard/trends/patches` | Patch trend |
| GET | `/api/v1/dashboard/trends/devices` | Device count trend |
| GET | `/api/v1/dashboard/distribution/os` | Devices by OS |
| GET | `/api/v1/dashboard/distribution/department` | Devices by department |
| GET | `/api/v1/dashboard/incidents` | Incident breakdown |
| GET | `/api/v1/dashboard/vulnerabilities` | Vulnerability breakdown |
| GET | `/api/v1/dashboard/sla` | SLA breach details |
| GET | `/api/v1/dashboard/hr-metrics` | HR metrics |
| GET | `/api/v1/dashboard/risk-heat-map` | Risk by org unit |
| GET | `/api/v1/dashboard/activity-feed` | Recent events |
| GET | `/api/v1/dashboard/export/csv` | CSV export |
| GET | `/api/v1/dashboard/export/pdf` | PDF report |

### 7.5 Sync & Connector APIs

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/sync/trigger` | Trigger sync job |
| GET | `/api/v1/sync/status` | Get sync job status |
| GET | `/api/v1/connectors` | List configured connectors |
| POST | `/api/v1/connectors` | Register connector |
| GET | `/api/v1/connectors/:id` | Get connector details |
| PUT | `/api/v1/connectors/:id` | Update connector |
| POST | `/api/v1/connectors/:id/test` | Test connectivity |
| POST | `/api/v1/connectors/:id/sync` | Trigger manual sync |
| GET | `/api/v1/connectors/:id/health` | Connector health |

### 7.6 Audit APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/audit/events` | Query audit log |
| GET | `/api/v1/audit/export` | Export audit log |
| GET | `/api/v1/audit/devices/:id` | Device audit trail |
| GET | `/api/v1/audit/users/:id` | User audit trail |

---

## 8. Connector Framework

### 8.1 Connector Interface

Every connector **must** implement the following interface:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONNECTOR INTERFACE (IConnector)                          │
│                                                                             │
│  LIFECYCLE                                                                  │
│  ├── initialize(config: ConnectorConfig) → void                            │
│  ├── connect() → ConnectionStatus                                          │
│  ├── disconnect() → void                                                   │
│  ├── health_check() → HealthStatus                                         │
│  └── get_capabilities() → ConnectorCapabilities                            │
│                                                                             │
│  ENTITY PROVIDERS                                                           │
│  ├── get_devices(cursor?, page_size?, filters?) → PaginatedResult          │
│  ├── get_device(device_id) → CanonicalDevice | null                        │
│  ├── get_users(cursor?, page_size?, filters?) → PaginatedResult            │
│  ├── get_applications(cursor?, page_size?, filters?) → PaginatedResult     │
│  └── get_patch_status(device_id?, cursor?, page_size?, filters?) → Result  │
│                                                                             │
│  ACTION PROVIDERS                                                           │
│  ├── deploy_application(app_id, device_ids, params?) → DeploymentResult   │
│  ├── run_script(device_id, content, type, timeout?, params?) → Result     │
│  └── restart_device(device_id, force?, reason?) → ActionResult            │
│                                                                             │
│  EXTENDED (optional)                                                        │
│  ├── lock_device(device_id) → ActionResult                                 │
│  ├── wipe_device(device_id) → ActionResult                                 │
│  ├── remote_shell(device_id) → ShellSession                               │
│  └── create_incident(title, description, priority) → Incident             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Connector Capability Matrix

| Capability | Ivanti Neurons | Microsoft Intune | Microsoft Entra ID | ServiceNow | N-Able |
|---|---|---|---|---|---|
| get_devices | ✅ Agent + unmanaged | ✅ Managed devices | ✅ Registered devices | ✅ CMDB CIs | ✅ Monitored |
| get_users | ✅ Identity directory | ✅ Graph users | ✅ Primary source | ✅ sys_user | ❌ |
| get_applications | ✅ Software catalog | ✅ Mobile apps | ✅ Service principals | ✅ Service CI | ✅ |
| get_patch_status | ✅ Patch compliance | ✅ Compliance policies | ❌ | ✅ Vuln integration | ✅ |
| deploy_application | ✅ Software deploy | ✅ App assignment | ⚠️ App role | ✅ Change request | ✅ |
| run_script | ✅ Agent script | ✅ PS/Shell script | ❌ | ✅ MID Server | ✅ Agent |
| restart_device | ✅ Agent command | ✅ RebootNow | ❌ | ✅ Orchestration | ✅ Agent |
| remote_shell | ✅ Agent shell | ❌ | ❌ | ✅ MID Server | ✅ Agent shell |
| lock_device | ❌ | ✅ | ❌ | ❌ | ❌ |
| wipe_device | ❌ | ✅ | ❌ | ❌ | ❌ |
| create_incident | ✅ | ❌ | ❌ | ✅ | ❌ |

### 8.3 Canonical Data Models

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CANONICAL ENTITY MODELS                                   │
│                                                                             │
│  CanonicalDevice:                                                           │
│  ├── id, external_id, connector_type, tenant_id                            │
│  ├── canonical_name, asset_type, os_type, os_version                       │
│  ├── serial_number, manufacturer, model                                    │
│  ├── primary_user_id, ip_address, mac_address                              │
│  ├── status, compliance_status, last_seen, metadata (JSONB)               │
│  └── raw_payload (JSONB)                                                   │
│                                                                             │
│  CanonicalUser:                                                             │
│  ├── id, external_id, connector_type, tenant_id                            │
│  ├── email, display_name, first_name, last_name                            │
│  ├── department, job_title, manager_id, status, roles[]                    │
│  └── raw_payload (JSONB)                                                   │
│                                                                             │
│  CanonicalApplication:                                                      │
│  ├── id, external_id, connector_type, tenant_id                            │
│  ├── name, version, publisher, category, cpe_identifier                   │
│  ├── is_approved, install_command, uninstall_command                       │
│  └── raw_payload (JSONB)                                                   │
│                                                                             │
│  CanonicalPatchStatus:                                                      │
│  ├── id, device_id, patch_name, kb_article, severity                      │
│  ├── status, release_date, install_date, is_superseded                     │
│  └── connector_type, metadata (JSONB)                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.4 Connector Runtime Features

| Feature | Implementation |
|---|---|
| Authentication | OAuth 2.0, API Key, Certificate — pluggable per connector |
| Rate Limiting | Token bucket algorithm, per-connector configuration |
| Circuit Breaker | Automatic fault isolation with recovery |
| Retry Policy | Exponential backoff with jitter |
| Pagination | Cursor-based for large datasets |
| Event Emission | Structured events to Kafka for all operations |
| Health Monitoring | Periodic health checks with alerting |
| Metrics | OpenTelemetry metrics for latency, errors, throughput |

---

## 9. Security Model

### 9.1 Authentication Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                                       │
│                                                                             │
│  User Browser          API Gateway          Identity Service    Entra ID   │
│       │                    │                      │                │        │
│       │  Login redirect    │                      │                │        │
│       │───────────────────>│                      │                │        │
│       │                    │  Validate with        │                │        │
│       │                    │  Entra ID            │                │        │
│       │                    │─────────────────────>│───────────────>│        │
│       │                    │                      │                │        │
│       │                    │                      │  SAML/OIDC     │        │
│       │                    │                      │  Assertion     │        │
│       │                    │                      │<───────────────│        │
│       │                    │                      │                │        │
│       │                    │  JWT token           │                │        │
│       │                    │<─────────────────────│                │        │
│       │  Redirect with     │                      │                │        │
│       │  session cookie    │                      │                │        │
│       │<───────────────────│                      │                │        │
│       │                    │                      │                │        │
│       │  API calls with    │                      │                │        │
│       │  Bearer token      │                      │                │        │
│       │───────────────────>│  Validate JWT        │                │        │
│       │                    │  (local, no          │                │        │
│       │                    │   round-trip)        │                │        │
│       │                    │                      │                │        │
│       │  Response           │                      │                │        │
│       │<───────────────────│                      │                │        │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 JWT Token Structure

```json
{
  "sub": "user-uuid",
  "email": "jane.smith@contoso.com",
  "name": "Jane Smith",
  "roles": ["operations_manager"],
  "tenant_id": "tenant-uuid",
  "org_ids": ["org-uuid-1", "org-uuid-2"],
  "permissions": ["device:read", "device:write", "action:execute", "action:approve"],
  "iat": 1718726400,
  "exp": 1718730000,
  "iss": "https://portal.contoso.com",
  "aud": "https://api.contoso.com"
}
```

### 9.3 Tenant Isolation

| Layer | Isolation Mechanism |
|---|---|
| Database | Row-Level Security (RLS) policies per tenant_id |
| API | JWT tenant claim validated by middleware; queries scoped automatically |
| Cache | Key prefixing: `device:{tenant_id}:{device_id}` |
| Search | Index aliasing per tenant in Elasticsearch |
| Connector | Per-tenant worker pools with isolated state |
| Network | Kubernetes NetworkPolicies per tenant namespace |

### 9.4 RBAC Policy Engine

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RBAC POLICY STRUCTURE                                     │
│                                                                             │
│  Role: Operations Manager                                                   │
│  ├── permissions:                                                           │
│  │   ├── resource: "device"                                                │
│  │   │   └── actions: ["read"]                                            │
│  │   ├── resource: "action"                                                │
│  │   │   └── actions: ["execute", "approve"]                              │
│  │   │   └── constraints:                                                 │
│  │   │       ├── action_type: ["run_script", "restart_device"]           │
│  │   │       ├── target_count_max: 50                                     │
│  │   │       └── risk_level_max: "high"                                  │
│  │   └── resource: "audit_log"                                            │
│  │       └── actions: ["read"]                                            │
│  ├── scope:                                                                 │
│  │   ├── type: "organization"                                             │
│  │   └── ids: ["org-uuid-engineering", "org-uuid-ops"]                   │
│  └── conditions:                                                            │
│      └── time_restrictions: "business_hours"  // Optional                  │
│                                                                             │
│  Evaluation: ALLOW if ANY role grants permission on the resource           │
│  with all constraints satisfied and scope includes the target.             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.5 Security Controls Summary

| Control | Implementation |
|---|---|
| Transport Security | TLS 1.3 everywhere; mTLS between services |
| Identity | Microsoft Entra ID with SAML 2.0 / OIDC |
| MFA | Enforced at Entra ID Conditional Access |
| Session Management | JWT (15-min expiry), Refresh tokens (8-hour), Redis-backed |
| API Security | Rate limiting (100 req/min per user), input validation, CORS |
| Data Protection | AES-256 at rest; column-level encryption for PII |
| Secrets | HashiCorp Vault with dynamic credentials and auto-rotation |
| Network | Private subnets, NetworkPolicies, WAF, DDoS protection |
| Logging | Immutable audit trail, SIEM integration ready |
| Vulnerability | SAST/DAST in CI/CD, dependency scanning, penetration testing |

---

## 10. Deployment Model

### 10.1 Infrastructure Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT ARCHITECTURE                                    │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  CLOUD REGION (Primary)                                               │  │
│  │                                                                       │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  KUBERNETE CLUSTER (EKS / AKS / GKE)                            │  │  │
│  │  │                                                                 │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │  │  │
│  │  │  │ Frontend     │  │ API Gateway  │  │ Core Services        │  │  │  │
│  │  │  │ Namespace    │  │ Namespace    │  │ Namespace            │  │  │  │
│  │  │  │ (3 replicas) │  │ (3 replicas) │  │ (3-5 replicas each)  │  │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │  │  │
│  │  │                                                                 │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │  │  │
│  │  │  │ Connector    │  │ Observability│  │ Security             │  │  │  │
│  │  │  │ Namespace    │  │ Namespace    │  │ Namespace            │  │  │  │
│  │  │  │ (2 per conn) │  │ (Grafana,    │  │ (Vault, OPA,        │  │  │  │
│  │  │  │              │  │  Loki, etc)  │  │  Cert Manager)       │  │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                       │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  DATA INFRASTRUCTURE                                            │  │  │
│  │  │                                                                 │  │  │
│  │  │  PostgreSQL HA    Redis Cluster    Kafka Cluster    Elasticsearch│  │  │
│  │  │  (Primary + 2RR)  (3+3)           (3 brokers)      (3 nodes)   │  │  │
│  │  │                                                                 │  │  │
│  │  │  ClickHouse       RabbitMQ         Vault Cluster                 │  │  │
│  │  │  (3 nodes)        (3 nodes)        (3 nodes)                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  DR REGION (Secondary)                                                │  │
│  │  • Async replication of PostgreSQL                                    │  │
│  │  • Kafka mirror maker                                                 │  │
│  │  • RTO: < 30 minutes, RPO: < 5 minutes                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Environment Strategy

| Environment | Purpose | Data | Scaling |
|---|---|---|---|
| **Development** | Feature development, unit tests | Synthetic seed data | 1 replica per service |
| **Staging** | Integration testing, QA | Anonymized production snapshot | 2 replicas per service |
| **Production** | Live system | Real data | Full HA, auto-scaling |
| **DR** | Disaster recovery | Async replicated | Minimal (scales on failover) |

### 10.3 CI/CD Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CI/CD PIPELINE                                             │
│                                                                             │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌────────┐ │
│  │  Code    │──>│  Build   │──>│  Test    │──>│  Stage   │──>│  Prod  │ │
│  │  Commit  │   │  & Scan  │   │  & QA    │   │  Deploy  │   │ Deploy │ │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘   └────────┘ │
│       │              │              │               │              │       │
│       ▼              ▼              ▼               ▼              ▼       │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌────────┐ │
│  │ GitHub   │   │ Docker   │   │ Unit     │   │ Staging  │   │ Blue/  │ │
│  │ Actions  │   │ Build    │   │ Integ.   │   │ K8s      │   │ Green  │ │
│  │ Trigger  │   │ SAST     │   │ E2E      │   │ Namespace│   │ Switch │ │
│  │          │   │ SCA      │   │ Contract │   │          │   │        │ │
│  │          │   │ Lint     │   │ Perf     │   │          │   │        │ │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘   └────────┘ │
│                                                                             │
│  Quality Gates:                                                             │
│  • SAST: No critical/high findings                                         │
│  • Unit tests: > 80% coverage                                              │
│  • Integration: All connector contracts pass                               │
│  • E2E: All critical user journeys pass                                    │
│  • Performance: p95 < target thresholds                                    │
│  • Security: OWASP scan clean                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.4 High Availability

| Component | HA Strategy | RPO | RTO |
|---|---|---|---|
| PostgreSQL | Streaming replication (sync) + Patroni failover | 0 | < 30s |
| Redis | Redis Sentinel / Cluster (3+3) | < 1s | < 10s |
| Kafka | Replication factor 3, min.insync.replicas=2 | 0 | < 30s |
| Elasticsearch | 3-node cluster, replica shards | 0 | < 30s |
| Kubernetes | Multi-AZ, pod anti-affinity, HPA | N/A | < 60s |
| **DR (Cross-Region)** | Async replication | < 5 min | < 30 min |

---

## 11. MVP Scope

### 11.1 MVP Definition

The MVP delivers core inventory and basic remote operations capabilities for a single tenant, establishing the platform foundation for future expansion.

### 11.2 MVP Inclusions

| Category | Feature | Phase |
|---|---|---|
| **Foundation** | Kubernetes cluster, CI/CD pipeline, observability stack | Week 1-2 |
| **Foundation** | Identity & RBAC service (Entra ID integration) | Week 1-2 |
| **Foundation** | Connector SDK & runtime framework | Week 1-2 |
| **Foundation** | PostgreSQL, Redis, Kafka, Elasticsearch setup | Week 1-2 |
| **Inventory** | Ivanti Neurons connector (devices, users) | Week 3-4 |
| **Inventory** | Microsoft Intune connector (devices, users) | Week 3-4 |
| **Inventory** | Microsoft Entra ID connector (devices, users) | Week 3-4 |
| **Inventory** | Merge engine with matching rules | Week 4-5 |
| **Inventory** | Unified device list API | Week 5 |
| **Inventory** | Basic device list UI with search and filters | Week 5-6 |
| **Operations** | Remote Ops Service — action routing engine | Week 6-7 |
| **Operations** | Script execution (single device) | Week 7 |
| **Operations** | Device restart (single device) | Week 7 |
| **Operations** | Action status tracking | Week 7-8 |
| **Operations** | Basic action console UI | Week 8 |
| **Dashboard** | Dashboard Service with summary API | Week 8-9 |
| **Dashboard** | Executive dashboard UI (5 KPIs: Total, Online, Offline, Compliance, Patch) | Week 9-10 |
| **Audit** | Audit logging pipeline | Week 9 |
| **Audit** | Basic audit log viewer | Week 10 |

### 11.3 MVP Exclusions

| Feature | Reason | Target Phase |
|---|---|---|
| ServiceNow integration | Phase 2 priority | Phase 2 |
| Bulk operations (>50 devices) | Scale feature | Phase 2 |
| Remote shell | Requires session management | Phase 2 |
| Remediation engine | Complex policy engine | Phase 3 |
| Patch management UI | Depends on full patch engine | Phase 2 |
| Advanced dashboards (risk heat map, HR) | Analytics features | Phase 3 |
| Multi-tenant | Platform feature | Phase 2 |
| N-Able / Addigy / Defender | Future integrations | Phase 4 |
| Power BI integration | Analytics | Phase 4 |
| PDF export | Reporting | Phase 3 |
| Mobile app | Future | Phase 4 |

### 11.4 MVP Success Criteria

| Criterion | Target |
|---|---|
| Ingest 10,000 devices from 3 sources | ✅ |
| Merge accuracy (auto-merge correct rate) | > 95% |
| Device list query response | < 100ms (p95) |
| Script execution success rate | > 98% |
| Dashboard load time | < 2 seconds |
| System uptime | > 99.5% |
| Zero critical security findings | ✅ |
| Full audit trail for all operations | ✅ |

### 11.5 MVP Timeline

```
Week:  1  2  3  4  5  6  7  8  9  10
       ├──────┤
       Foundation
              ├────────┤
              Connectors (Ivanti, Intune, Entra)
                       ├────┤
                       Merge Engine
                            ├────┤
                            Inventory API + UI
                                  ├────┤
                                  Remote Ops + Console
                                       ├────┤
                                       Dashboard + Audit
                                            ┤
                                            UAT & Launch
```

---

## 12. Future Scope

### 12.1 Phase 2 — Core Capabilities (Months 4-6)

| Feature | Description |
|---|---|
| ServiceNow integration | CMDB sync, incident creation, change requests |
| Multi-tenant platform | Tenant management, namespace isolation, tenant-aware billing |
| Bulk operations | Actions on up to 1,000 devices with parallel execution |
| Remote shell | Interactive shell sessions with recording |
| Patch management | Full patch policy engine, compliance dashboard, deployment scheduling |
| Application deployment | App catalog, approval workflow, selective deployment |
| Advanced search | Elasticsearch-powered full-text search with facets |
| Dashboard expansion | All 10 KPIs, trend charts, risk heat map |
| Export capabilities | CSV and PDF report generation |
| Notification service | Email, Slack, Teams notifications for critical events |

### 12.2 Phase 3 — Intelligence & Scale (Months 7-9)

| Feature | Description |
|---|---|
| Remediation engine | Policy-based automated remediation with approval workflows |
| ML-based anomaly detection | Unusual device behavior, security threat detection |
| Predictive patching | Risk-based patch prioritization using ML |
| Self-service portal | End users can view own devices, request actions, track status |
| Advanced reporting | Custom report builder with scheduling |
| Compliance automation | Automated compliance remediation workflows |
| Workflow templates | Reusable action sequences with conditional logic |
| API marketplace | Third-party connector SDK and marketplace |

### 12.3 Phase 4 — Future Integrations (Months 10-12)

| Feature | Description |
|---|---|
| Microsoft Defender connector | Threat signals, vulnerability assessment, incident correlation |
| N-Able connector | RMM data, monitoring alerts, automated remediation |
| Addigy connector | Apple device management, macOS/iOS-specific policies |
| Power BI integration | Executive dashboards, advanced analytics, custom reports |
| ServiceNow bidirectional sync | Full CMDB reconciliation, incident auto-creation |
| ITSM integration | Change management integration, approval routing |
| Mobile application | iOS/Android app for on-the-go operations |
| ChatOps integration | Slack/Teams bot for action execution |

### 12.4 Long-Term Vision (Year 2+)

| Feature | Description |
|---|---|
| AIOps | Intelligent root cause analysis, predictive maintenance |
| Zero-touch provisioning | Automated device onboarding with policy assignment |
| IoT device management | Extended device types (sensors, cameras, building systems) |
| Edge computing management | Manage edge devices and containers |
| Global multi-region | Active-active deployment across regions |
| FedRAMP/IL5 compliance | Government cloud readiness |
| Digital experience monitoring | End-user experience scoring, proactive issue detection |

---

## 13. Appendices

### Appendix A: Glossary

| Term | Definition |
|---|---|
| Canonical Model | Normalized data representation independent of source system |
| Connector | Adapter that maps a vendor system's API to the canonical interface |
| Merge Engine | Component that deduplicates and merges records from multiple sources |
| Confidence Score | Numeric value (0.0–1.0) indicating merge accuracy |
| Tenant | Isolated organizational unit with its own data and configuration |
| RBAC | Role-Based Access Control |
| RLS | Row-Level Security (PostgreSQL feature) |
| Temporal | Workflow engine for durable execution |
| Circuit Breaker | Pattern that prevents cascade failures |
| Materialized View | Pre-computed query result stored for fast access |

### Appendix B: Reference Architecture Diagrams

All diagrams referenced in this document are maintained in the project's architecture repository at:
- `/docs/architecture/system-architecture.md`
- `/docs/architecture/database-er-diagram.md`
- `/docs/architecture/api-flows.md`
- `/docs/architecture/connector-framework.md`
- `/docs/architecture/deployment-architecture.md`
- `/docs/architecture/security-architecture.md`

### Appendix C: API Contract Specifications

Full OpenAPI 3.0 specifications maintained at:
- `/api/specs/inventory.yaml`
- `/api/specs/operations.yaml`
- `/api/specs/dashboard.yaml`
- `/api/specs/audit.yaml`
- `/api/specs/connectors.yaml`

### Appendix D: Revision History

| Version | Date | Author | Changes |
|---|---|---|---|
| 0.1 | 2026-06-10 | Enterprise Architecture | Initial draft |
| 0.2 | 2026-06-15 | Enterprise Architecture | Added connector framework, security model |
| 0.3 | 2026-06-17 | Enterprise Architecture | Added operations center, dashboard specs |
| 1.0 | 2026-06-18 | Enterprise Architecture | Baseline release — consolidated all prior designs |

---

**End of Document**

*This document is the source of truth for all development activities. Any deviation requires written approval from the Enterprise Architecture team.*
