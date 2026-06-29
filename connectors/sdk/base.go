package sdk

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// --- Canonical Data Models ---

type CanonicalDevice struct {
	ID               string                 `json:"id"`
	ExternalID       string                 `json:"external_id"`
	ConnectorType    string                 `json:"connector_type"`
	TenantID         string                 `json:"tenant_id"`
	CanonicalName    string                 `json:"canonical_name"`
	AssetType        string                 `json:"asset_type"`
	OSType           string                 `json:"os_type"`
	OSVersion        string                 `json:"os_version"`
	SerialNumber     string                 `json:"serial_number"`
	Manufacturer     string                 `json:"manufacturer"`
	Model            string                 `json:"model"`
	PrimaryUserID    string                 `json:"primary_user_id"`
	IPAddress        string                 `json:"ip_address"`
	MACAddress       string                 `json:"mac_address"`
	Status           string                 `json:"status"`
	ComplianceStatus string                 `json:"compliance_status"`
	LastSeen         time.Time              `json:"last_seen"`
	RiskScore        int                    `json:"risk_score"`
	Metadata         map[string]interface{} `json:"metadata"`
}

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

// --- Connector Interface ---

type Connector interface {
	Initialize(config ConnectorConfig) error
	Connect() error
	Disconnect() error
	HealthCheck() HealthStatus
	GetCapabilities() ConnectorCapabilities
	GetDevices(cursor string, pageSize int, filters map[string]interface{}) ([]CanonicalDevice, string, int, error)
	GetDevice(deviceID string) (*CanonicalDevice, error)
	GetUsers(cursor string, pageSize int, filters map[string]interface{}) ([]CanonicalUser, string, int, error)
	RunScript(deviceID string, content string, scriptType string, timeout int) (ScriptResult, error)
	RestartDevice(deviceID string, force bool, reason string) (ActionResult, error)
	DeployApplication(appID string, deviceIDs []string) (DeploymentResult, error)
}

// --- Supporting Types ---

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

type HealthStatus struct {
	ConnectorType      string `json:"connector_type"`
	Status             string `json:"status"`
	LastSuccessfulSync string `json:"last_successful_sync,omitempty"`
	ErrorMessage       string `json:"error_message,omitempty"`
	LatencyMs          int    `json:"latency_ms"`
}

type ScriptResult struct {
	ScriptID string  `json:"script_id"`
	DeviceID string  `json:"device_id"`
	ExitCode int     `json:"exit_code"`
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Duration float64 `json:"duration_seconds"`
	Status   string  `json:"status"`
}

type ActionResult struct {
	ActionID   string `json:"action_id"`
	ActionType string `json:"action_type"`
	DeviceID   string `json:"device_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type DeploymentResult struct {
	DeploymentID string `json:"deployment_id"`
	DeviceID     string `json:"device_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

type ConnectorCapabilities struct {
	ConnectorType string   `json:"connector_type"`
	EntityTypes   []string `json:"entity_types"`
	Operations    []string `json:"operations"`
	Actions       []string `json:"actions"`
	AuthMethods   []string `json:"auth_methods"`
	SyncModes     []string `json:"sync_modes"`
	RateLimit     int      `json:"rate_limit"`
}

// --- Base Connector ---

type BaseConnector struct {
	Config      ConnectorConfig
	mu          sync.RWMutex
	connected   bool
	healthy     bool
	rateLimiter *TokenBucket
}

type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	rate       float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(rate, burst int) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(burst),
		maxTokens:  float64(burst),
		rate:       float64(rate) / 60.0,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Acquire(ctx context.Context) error {
	for {
		tb.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(tb.lastRefill).Seconds()
		tb.tokens = min(tb.maxTokens, tb.tokens+elapsed*tb.rate)
		tb.lastRefill = now

		if tb.tokens >= 1 {
			tb.tokens--
			tb.mu.Unlock()
			return nil
		}
		tb.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (b *BaseConnector) BaseInitialize(config ConnectorConfig, name string) {
	b.Config = config
	b.rateLimiter = NewTokenBucket(config.RateLimit, config.RateLimit*2)
}

func (b *BaseConnector) BaseConnect() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connected = true
	b.healthy = true
}

func (b *BaseConnector) BaseDisconnect() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connected = false
	b.healthy = false
}

func (b *BaseConnector) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.connected
}

// --- Registry ---

type Registry struct {
	mu         sync.RWMutex
	connectors map[string]func() Connector
}

func NewRegistry() *Registry {
	return &Registry{
		connectors: make(map[string]func() Connector),
	}
}

func (r *Registry) Register(connectorType string, factory func() Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connectors[connectorType] = factory
}

func (r *Registry) Create(connectorType string) (Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.connectors[connectorType]
	if !ok {
		return nil, fmt.Errorf("unknown connector type: %s", connectorType)
	}
	return factory(), nil
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.connectors))
	for t := range r.connectors {
		types = append(types, t)
	}
	return types
}

var GlobalRegistry = NewRegistry()

func RegisterConnector(connectorType string, factory func() Connector) {
	GlobalRegistry.Register(connectorType, factory)
}
