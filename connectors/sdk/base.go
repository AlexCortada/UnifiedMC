package sdk

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Connector is the interface that every connector must implement
type Connector interface {
	Initialize(config ConnectorConfig) error
	Connect() error
	Disconnect() error
	HealthCheck() HealthStatus
	GetCapabilities() ConnectorCapabilities
	GetDevices(cursor string, pageSize int, filters map[string]interface{}) ([]types.CanonicalDevice, string, int, error)
	GetDevice(deviceID string) (*types.CanonicalDevice, error)
	GetUsers(cursor string, pageSize int, filters map[string]interface{}) ([]types.CanonicalUser, string, int, error)
	RunScript(deviceID string, content string, scriptType string, timeout int) (types.ScriptResult, error)
	RestartDevice(deviceID string, force bool, reason string) (types.ActionResult, error)
	DeployApplication(appID string, deviceIDs []string) (types.DeploymentResult, error)
}

// BaseConnector provides shared infrastructure for all connectors
type BaseConnector struct {
	Config      ConnectorConfig
	Logger      *log.Logger
	mu          sync.RWMutex
	connected   bool
	healthy     bool
	rateLimiter *TokenBucket
}

// TokenBucket implements a simple rate limiter
type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	rate       float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate, burst int) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(burst),
		maxTokens:  float64(burst),
		rate:       float64(rate) / 60.0,
		lastRefill: time.Now(),
	}
}

// Acquire blocks until a token is available
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

// BaseInitialize sets up the base connector
func (b *BaseConnector) BaseInitialize(config ConnectorConfig, name string) {
	b.Config = config
	b.Logger = log.Printf("[%s] ", name)
	b.rateLimiter = NewTokenBucket(config.RateLimit, config.RateLimit*2)
}

// BaseConnect marks the connector as connected
func (b *BaseConnector) BaseConnect() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connected = true
	b.healthy = true
}

// BaseDisconnect marks the connector as disconnected
func (b *BaseConnector) BaseDisconnect() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connected = false
	b.healthy = false
}

// IsConnected returns connection status
func (b *BaseConnector) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.connected
}

// Registry manages all connector implementations
type Registry struct {
	mu         sync.RWMutex
	connectors map[string]func() Connector
}

// NewRegistry creates a new connector registry
func NewRegistry() *Registry {
	return &Registry{
		connectors: make(map[string]func() Connector),
	}
}

// Register registers a connector factory function
func (r *Registry) Register(connectorType string, factory func() Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connectors[connectorType] = factory
}

// Create creates a new connector instance by type
func (r *Registry) Create(connectorType string) (Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.connectors[connectorType]
	if !ok {
		return nil, fmt.Errorf("unknown connector type: %s", connectorType)
	}
	return factory(), nil
}

// List returns all registered connector types
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.connectors))
	for t := range r.connectors {
		types = append(types, t)
	}
	return types
}

// Global registry instance
var GlobalRegistry = NewRegistry()

// RegisterConnector registers a connector in the global registry
func RegisterConnector(connectorType string, factory func() Connector) {
	GlobalRegistry.Register(connectorType, factory)
}

// Ensure types package is imported
var _ = types.CanonicalDevice{}
