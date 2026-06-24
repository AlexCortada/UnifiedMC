package types

import (
	"time"
)

// AssetType represents the type of device
type AssetType string

const (
	Workstation AssetType = "workstation"
	Server      AssetType = "server"
	Mobile      AssetType = "mobile"
	IoT         AssetType = "iot"
	VM          AssetType = "vm"
	Container   AssetType = "container"
	Unknown     AssetType = "unknown"
)

// OSType represents the operating system
type OSType string

const (
	Windows OSType = "windows"
 MacOS   OSType = "macos"
	Linux   OSType = "linux"
	iOS     OSType = "ios"
	Android OSType = "android"
	Other   OSType = "other"
)

// AssetStatus represents device status
type AssetStatus string

const (
	Active       AssetStatus = "active"
	Inactive     AssetStatus = "inactive"
	Retired      AssetStatus = "retired"
	Quarantined  AssetStatus = "quarantined"
	Provisioning AssetStatus = "provisioning"
)

// ComplianceStatus represents compliance state
type ComplianceStatus string

const (
	Compliant    ComplianceStatus = "compliant"
	NonCompliant ComplianceStatus = "non_compliant"
	Unknown      ComplianceStatus = "unknown"
	Exempt       ComplianceStatus = "exempt"
)

// PatchStatus represents patch installation state
type PatchStatus string

const (
	Installed  PatchStatus = "installed"
	Missing    PatchStatus = "missing"
	Pending    PatchStatus = "pending"
	Failed     PatchStatus = "failed"
	Superseded PatchStatus = "superseded"
)

// Severity represents vulnerability severity
type Severity string

const (
	Critical Severity = "critical"
	High     Severity = "high"
	Medium   Severity = "medium"
	Low      Severity = "low"
)

// ConnectionStatus represents connector connection state
type ConnectionStatus string

const (
	Connected   ConnectionStatus = "connected"
	Disconnected ConnectionStatus = "disconnected"
	Degraded    ConnectionStatus = "degraded"
	AuthFailed  ConnectionStatus = "auth_failed"
)

// ActionType represents the type of remote action
type ActionType string

const (
	RunScript       ActionType = "run_script"
	RestartDevice   ActionType = "restart_device"
	DeployApplication ActionType = "deploy_application"
	PatchDevice     ActionType = "patch_device"
	RemoteShell     ActionType = "remote_shell"
)

// ActionStatus represents the state of an action
type ActionStatus string

const (
	Pending           ActionStatus = "pending"
	Validating        ActionStatus = "validating"
	AwaitingApproval ActionStatus = "awaiting_approval"
	Approved          ActionStatus = "approved"
	Queued            ActionStatus = "queued"
	Executing        ActionStatus = "executing"
	Completed         ActionStatus = "completed"
	Failed            ActionStatus = "failed"
	Cancelled         ActionStatus = "cancelled"
	Partial           ActionStatus = "partial"
)

// CanonicalDevice is the normalized device representation across all connectors
type CanonicalDevice struct {
	ID                string                 `json:"id"`
	ExternalID        string                 `json:"external_id"`
	ConnectorType     string                 `json:"connector_type"`
	TenantID          string                 `json:"tenant_id"`
	CanonicalName     string                 `json:"canonical_name"`
	AssetType         AssetType              `json:"asset_type"`
	OSType            OSType                 `json:"os_type"`
	OSVersion         string                 `json:"os_version"`
	SerialNumber      string                 `json:"serial_number"`
	Manufacturer      string                 `json:"manufacturer"`
	Model             string                 `json:"model"`
	PrimaryUserID     string                 `json:"primary_user_id"`
	IPAddress         string                 `json:"ip_address"`
	MACAddress        string                 `json:"mac_address"`
	Status            AssetStatus            `json:"status"`
	ComplianceStatus  ComplianceStatus       `json:"compliance_status"`
	LastSeen          time.Time              `json:"last_seen"`
	RiskScore         int                    `json:"risk_score"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// CanonicalUser is the normalized user representation
type CanonicalUser struct {
	ID            string   `json:"id"`
	ExternalID    string   `json:"external_id"`
	ConnectorType string   `json:"connector_type"`
	TenantID      string   `json:"tenant_id"`
	Email         string   `json:"email"`
	DisplayName   string   `json:"display_name"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	Department    string   `json:"department"`
	JobTitle      string   `json:"job_title"`
	ManagerID     string   `json:"manager_id"`
	Status        string   `json:"status"`
	Roles         []string `json:"roles"`
	Groups        []string `json:"groups"`
}

// CanonicalApplication is the normalized application representation
type CanonicalApplication struct {
	ID            string                 `json:"id"`
	ExternalID    string                 `json:"external_id"`
	ConnectorType string                 `json:"connector_type"`
	TenantID      string                 `json:"tenant_id"`
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Publisher     string                 `json:"publisher"`
	Category      string                 `json:"category"`
	CPEIdentifier string                 `json:"cpe_identifier"`
	IsApproved    bool                   `json:"is_approved"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// CanonicalPatchStatus is the normalized patch status
type CanonicalPatchStatus struct {
	ID             string       `json:"id"`
	DeviceID       string       `json:"device_id"`
	PatchName      string       `json:"patch_name"`
	KBArticle      string       `json:"kb_article"`
	Severity       Severity     `json:"severity"`
	Status         PatchStatus  `json:"status"`
	ReleaseDate    time.Time    `json:"release_date"`
	InstallDate    time.Time    `json:"install_date"`
	IsSuperseded   bool         `json:"is_superseded"`
	ConnectorType  string       `json:"connector_type"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ConnectorConfig is the initialization configuration for a connector
type ConnectorConfig struct {
	ConnectorType  string                 `json:"connector_type"`
	TenantID       string                 `json:"tenant_id"`
	Name           string                 `json:"name"`
	Enabled        bool                   `json:"enabled"`
	BaseURL        string                 `json:"base_url"`
	Auth           map[string]interface{} `json:"auth"`
	RateLimit      int                    `json:"rate_limit"`
	TimeoutSeconds int                    `json:"timeout_seconds"`
	MaxRetries     int                    `json:"max_retries"`
	SyncInterval   int                    `json:"sync_interval_minutes"`
	Metadata       map[string]interface{} `json:"metadata"`
	Filters        map[string]interface{} `json:"filters"`
}

// HealthStatus reports connector health
type HealthStatus struct {
	ConnectorType     string            `json:"connector_type"`
	Status            ConnectionStatus  `json:"status"`
	LastSuccessfulSync time.Time        `json:"last_successful_sync"`
	ErrorMessage      string            `json:"error_message"`
	LatencyMs         int               `json:"latency_ms"`
	Details           map[string]interface{} `json:"details"`
}

// ScriptResult is the output of script execution
type ScriptResult struct {
	ScriptID string  `json:"script_id"`
	DeviceID string  `json:"device_id"`
	ExitCode int     `json:"exit_code"`
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Duration float64 `json:"duration_seconds"`
	Status   string  `json:"status"`
}

// ActionResult is the output of a device action
type ActionResult struct {
	ActionID   string `json:"action_id"`
	ActionType string `json:"action_type"`
	DeviceID   string `json:"device_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// DeploymentResult is the output of app deployment
type DeploymentResult struct {
	DeploymentID string `json:"deployment_id"`
	DeviceID     string `json:"device_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// ConnectorCapabilities declares what a connector supports
type ConnectorCapabilities struct {
	ConnectorType string   `json:"connector_type"`
	EntityTypes   []string `json:"entity_types"`
	Operations    []string `json:"operations"`
	Actions       []string `json:"actions"`
	AuthMethods   []string `json:"auth_methods"`
	SyncModes     []string `json:"sync_modes"`
	RateLimit     int      `json:"rate_limit"`
}

// PaginatedResult wraps paginated responses
type PaginatedResult[T any] struct {
	Items       []T    `json:"items"`
	TotalCount  int    `json:"total_count"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	HasNext     bool   `json:"has_next"`
	NextCursor  string `json:"next_cursor"`
}
