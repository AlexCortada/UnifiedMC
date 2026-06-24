-- Unified IT Operations Portal - Initial Schema
-- PostgreSQL 16+

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =====================================================
-- TENANTS & ORGANIZATIONS
-- =====================================================

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    parent_org_id UUID REFERENCES organizations(id),
    path LTREE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organizations_tenant ON organizations(tenant_id);
CREATE INDEX idx_organizations_path ON organizations USING GIST(path);

-- =====================================================
-- USERS
-- =====================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    roles TEXT[] DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);

-- =====================================================
-- UNIFIED DEVICES (Canonical)
-- =====================================================

CREATE TABLE unified_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    org_id UUID REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    asset_type VARCHAR(50) NOT NULL DEFAULT 'unknown',
    os_type VARCHAR(50),
    os_version VARCHAR(100),
    serial_number VARCHAR(255),
    manufacturer VARCHAR(255),
    model VARCHAR(255),
    primary_user_id UUID REFERENCES users(id),
    department VARCHAR(255),
    compliance_status VARCHAR(50) DEFAULT 'unknown',
    patch_status_summary JSONB DEFAULT '{}',
    last_seen TIMESTAMPTZ,
    risk_score SMALLINT DEFAULT 50,
    installed_apps_count INT DEFAULT 0,
    merged_from_sources VARCHAR(50)[] DEFAULT '{}',
    merge_confidence FLOAT DEFAULT 0.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by_sync_at TIMESTAMPTZ
);

CREATE INDEX idx_unified_devices_tenant ON unified_devices(tenant_id);
CREATE INDEX idx_unified_devices_user ON unified_devices(primary_user_id);
CREATE INDEX idx_unified_devices_compliance ON unified_devices(compliance_status);
CREATE INDEX idx_unified_devices_last_seen ON unified_devices(last_seen DESC);
CREATE INDEX idx_unified_devices_sources ON unified_devices USING GIN(merged_from_sources);
CREATE INDEX idx_unified_devices_name_trgm ON unified_devices USING GIN(display_name gin_trgm_ops);

-- =====================================================
-- DEVICE SOURCES (Raw from connectors)
-- =====================================================

CREATE TABLE device_sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unified_device_id UUID REFERENCES unified_devices(id) ON DELETE SET NULL,
    connector_type VARCHAR(50) NOT NULL,
    external_id VARCHAR(500) NOT NULL,
    raw_data JSONB DEFAULT '{}',
    mapped_data JSONB DEFAULT '{}',
    last_sync_at TIMESTAMPTZ,
    sync_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    trust_score SMALLINT DEFAULT 50,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(connector_type, external_id)
);

CREATE INDEX idx_device_sources_unified ON device_sources(unified_device_id);
CREATE INDEX idx_device_sources_connector_ext ON device_sources(connector_type, external_id);
CREATE INDEX idx_device_sources_sync ON device_sources(sync_status, last_sync_at);

-- =====================================================
-- DEVICE SOURCE MAPPINGS
-- =====================================================

CREATE TABLE device_source_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unified_device_id UUID NOT NULL REFERENCES unified_devices(id) ON DELETE CASCADE,
    connector_type VARCHAR(50) NOT NULL,
    external_id VARCHAR(500) NOT NULL,
    match_confidence FLOAT NOT NULL,
    match_method VARCHAR(100) NOT NULL,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    matched_by VARCHAR(100) NOT NULL DEFAULT 'system'
);

CREATE INDEX idx_source_mappings_unified ON device_source_mappings(unified_device_id);

-- =====================================================
-- DEVICE APPLICATIONS
-- =====================================================

CREATE TABLE device_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unified_device_id UUID NOT NULL REFERENCES unified_devices(id) ON DELETE CASCADE,
    application_name VARCHAR(255) NOT NULL,
    application_version VARCHAR(100),
    publisher VARCHAR(255),
    install_date DATE,
    install_path VARCHAR(500),
    source_connector VARCHAR(50),
    is_approved BOOLEAN DEFAULT FALSE,
    raw_data JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_apps_device ON device_applications(unified_device_id);

-- =====================================================
-- DEVICE PATCHES
-- =====================================================

CREATE TABLE device_patches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unified_device_id UUID NOT NULL REFERENCES unified_devices(id) ON DELETE CASCADE,
    patch_name VARCHAR(255) NOT NULL,
    kb_article VARCHAR(50),
    severity VARCHAR(20) NOT NULL DEFAULT 'medium',
    status VARCHAR(20) NOT NULL DEFAULT 'missing',
    source_connector VARCHAR(50),
    detected_at TIMESTAMPTZ,
    installed_at TIMESTAMPTZ,
    superseded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_patches_device ON device_patches(unified_device_id);
CREATE INDEX idx_device_patches_severity ON device_patches(severity);

-- =====================================================
-- ACTION REQUESTS
-- =====================================================

CREATE TABLE action_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'normal',
    target_type VARCHAR(50),
    target_ids UUID[] NOT NULL,
    target_filter JSONB DEFAULT '{}',
    parameters JSONB DEFAULT '{}',
    backend_preference VARCHAR(50),
    selected_backend VARCHAR(50),
    fallback_backend VARCHAR(50),
    initiated_by UUID NOT NULL REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    requires_approval BOOLEAN DEFAULT FALSE,
    workflow_id UUID,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_seconds INT,
    result_summary JSONB DEFAULT '{}',
    error_details JSONB DEFAULT '{}',
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_action_requests_tenant ON action_requests(tenant_id);
CREATE INDEX idx_action_requests_status ON action_requests(status);

CREATE TABLE action_targets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action_request_id UUID NOT NULL REFERENCES action_requests(id) ON DELETE CASCADE,
    unified_device_id UUID NOT NULL REFERENCES unified_devices(id) ON DELETE CASCADE,
    backend_type VARCHAR(50) NOT NULL,
    external_id VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    result_code INT,
    output TEXT,
    error TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_seconds INT,
    retry_count INT DEFAULT 0,
    raw_result JSONB DEFAULT '{}'
);

CREATE INDEX idx_action_targets_request ON action_targets(action_request_id);
CREATE INDEX idx_action_targets_device ON action_targets(unified_device_id);

-- =====================================================
-- AUDIT LOG
-- =====================================================

CREATE TABLE action_audit_log (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    action_request_id UUID,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB DEFAULT '{}',
    actor_id UUID,
    actor_type VARCHAR(50) DEFAULT 'user',
    ip_address INET,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_tenant ON action_audit_log(tenant_id);
CREATE INDEX idx_audit_log_timestamp ON action_audit_log(timestamp DESC);

-- =====================================================
-- INCIDENTS
-- =====================================================

CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    incident_number VARCHAR(50) NOT NULL,
    external_id VARCHAR(500),
    source_connector VARCHAR(50) NOT NULL DEFAULT 'servicenow',
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority VARCHAR(10),
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    category VARCHAR(100),
    assignee_id UUID REFERENCES users(id),
    department VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    sla_response_deadline TIMESTAMPTZ,
    sla_resolution_deadline TIMESTAMPTZ,
    sla_response_breached BOOLEAN DEFAULT FALSE,
    sla_resolution_breached BOOLEAN DEFAULT FALSE,
    response_time_minutes INT,
    resolution_time_minutes INT,
    raw_data JSONB DEFAULT '{}'
);

CREATE INDEX idx_incidents_tenant ON incidents(tenant_id);
CREATE INDEX idx_incidents_status ON incidents(status);

-- =====================================================
-- CONNECTOR CONFIGURATIONS
-- =====================================================

CREATE TABLE connectors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    health VARCHAR(50) DEFAULT 'unknown',
    last_health_check TIMESTAMPTZ,
    last_successful_sync TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_connectors_tenant ON connectors(tenant_id);

-- =====================================================
-- MATERIALIZED VIEWS
-- =====================================================

CREATE MATERIALIZED VIEW mv_device_availability AS
SELECT
    tenant_id,
    COUNT(*) AS total_devices,
    COUNT(*) FILTER (WHERE last_seen > NOW() - INTERVAL '1 hour') AS online_devices,
    COUNT(*) FILTER (WHERE last_seen <= NOW() - INTERVAL '1 hour') AS offline_devices
FROM unified_devices
GROUP BY tenant_id;

CREATE MATERIALIZED VIEW mv_compliance_summary AS
SELECT
    tenant_id,
    COUNT(*) AS total_devices,
    COUNT(*) FILTER (WHERE compliance_status = 'compliant') AS compliant_count,
    COUNT(*) FILTER (WHERE compliance_status = 'non_compliant') AS non_compliant_count,
    ROUND(
        COUNT(*) FILTER (WHERE compliance_status = 'compliant')::NUMERIC /
        NULLIF(COUNT(*), 0) * 100, 1
    ) AS compliance_rate
FROM unified_devices
GROUP BY tenant_id;

