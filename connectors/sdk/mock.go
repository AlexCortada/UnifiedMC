package sdk

import (
	"context"
	"fmt"
	"time"
)

// MockConnector is a test implementation of the Connector interface.
type MockConnector struct {
	BaseConnector
	devices []CanonicalDevice
	users   []CanonicalUser
}

// NewMockConnector creates a new mock connector with sample data.
func NewMockConnector() *MockConnector {
	return &MockConnector{
		devices: generateMockDevices(),
		users:   generateMockUsers(),
	}
}

// Initialize prepares the mock connector.
func (m *MockConnector) Initialize(config ConnectorConfig) error {
	m.BaseInitialize(config, "mock")
	return nil
}

// Connect simulates connection.
func (m *MockConnector) Connect() error {
	m.BaseConnect()
	return nil
}

// Disconnect simulates disconnection.
func (m *MockConnector) Disconnect() error {
	m.BaseDisconnect()
	return nil
}

// HealthCheck returns mock health status.
func (m *MockConnector) HealthCheck() HealthStatus {
	return HealthStatus{
		ConnectorType: "mock",
		Status:        "connected",
		LatencyMs:     5,
	}
}

// GetCapabilities declares mock capabilities.
func (m *MockConnector) GetCapabilities() ConnectorCapabilities {
	return ConnectorCapabilities{
		ConnectorType: "mock",
		EntityTypes:   []string{"device", "user", "application"},
		Operations:    []string{"read", "write", "action"},
		Actions:       []string{"run_script", "restart_device", "deploy_application"},
		AuthMethods:   []string{"none"},
		SyncModes:     []string{"pull"},
		RateLimit:     1000,
	}
}

// GetDevices returns paginated mock devices.
func (m *MockConnector) GetDevices(cursor string, pageSize int, filters map[string]interface{}) ([]CanonicalDevice, string, int, error) {
	start := 0
	end := start + pageSize
	if end > len(m.devices) {
		end = len(m.devices)
	}
	nextCursor := ""
	if end < len(m.devices) {
		nextCursor = fmt.Sprintf("page-%d", start+pageSize)
	}
	return m.devices[start:end], nextCursor, len(m.devices), nil
}

// GetDevice returns a single mock device.
func (m *MockConnector) GetDevice(deviceID string) (*CanonicalDevice, error) {
	for _, d := range m.devices {
		if d.ExternalID == deviceID {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("device not found: %s", deviceID)
}

// GetUsers returns mock users.
func (m *MockConnector) GetUsers(cursor string, pageSize int, filters map[string]interface{}) ([]CanonicalUser, string, int, error) {
	return m.users, "", len(m.users), nil
}

// RunScript simulates script execution.
func (m *MockConnector) RunScript(deviceID string, content string, scriptType string, timeout int) (ScriptResult, error) {
	return ScriptResult{
		ScriptID: "mock-script",
		DeviceID: deviceID,
		ExitCode: 0,
		Stdout:   fmt.Sprintf("Mock script executed successfully on %s", deviceID),
		Duration: 1.5,
		Status:   "completed",
	}, nil
}

// RestartDevice simulates device restart.
func (m *MockConnector) RestartDevice(deviceID string, force bool, reason string) (ActionResult, error) {
	return ActionResult{
		ActionID:   fmt.Sprintf("restart-%d", time.Now().Unix()),
		ActionType: "restart",
		DeviceID:   deviceID,
		Status:     "success",
		Message:    "Restart command sent",
	}, nil
}

// DeployApplication simulates app deployment.
func (m *MockConnector) DeployApplication(appID string, deviceIDs []string) (DeploymentResult, error) {
	return DeploymentResult{
		DeploymentID: fmt.Sprintf("deploy-%d", time.Now().Unix()),
		DeviceID:     deviceIDs[0],
		Status:       "pending",
		Message:      fmt.Sprintf("Deployment queued for %d devices", len(deviceIDs)),
	}, nil
}

func generateMockDevices() []CanonicalDevice {
	return []CanonicalDevice{
		{
			ID:               "dev-001",
			ExternalID:       "mock-001",
			ConnectorType:    "mock",
			TenantID:         "tenant-001",
			CanonicalName:    "DESKTOP-ABC123",
			AssetType:        "workstation",
			OSType:           "windows",
			OSVersion:        "11 23H2",
			SerialNumber:     "SN-ABC-12345",
			Manufacturer:     "Dell",
			Model:            "OptiPlex 7090",
			PrimaryUserID:    "user-001",
			IPAddress:        "10.0.1.101",
			MACAddress:       "AA:BB:CC:DD:EE:01",
			Status:           "active",
			ComplianceStatus: "compliant",
			LastSeen:         time.Now().UTC().Add(-5 * time.Minute),
			Metadata:         map[string]interface{}{},
		},
		{
			ID:               "dev-002",
			ExternalID:       "mock-002",
			ConnectorType:    "mock",
			TenantID:         "tenant-001",
			CanonicalName:    "LAPTOP-XYZ789",
			AssetType:        "workstation",
			OSType:           "macos",
			OSVersion:        "14.2.1",
			SerialNumber:     "SN-XYZ-67890",
			Manufacturer:     "Apple",
			Model:            "MacBook Pro 16",
			PrimaryUserID:    "user-002",
			IPAddress:        "10.0.1.102",
			MACAddress:       "AA:BB:CC:DD:EE:02",
			Status:           "active",
			ComplianceStatus: "compliant",
			LastSeen:         time.Now().UTC().Add(-12 * time.Minute),
			Metadata:         map[string]interface{}{},
		},
		{
			ID:               "dev-003",
			ExternalID:       "mock-003",
			ConnectorType:    "mock",
			TenantID:         "tenant-001",
			CanonicalName:    "SERVER-WEB01",
			AssetType:        "server",
			OSType:           "linux",
			OSVersion:        "Ubuntu 22.04 LTS",
			SerialNumber:     "SN-SRV-11111",
			Manufacturer:     "HP",
			Model:            "ProLiant DL380",
			PrimaryUserID:    "user-003",
			IPAddress:        "10.0.1.201",
			MACAddress:       "AA:BB:CC:DD:EE:03",
			Status:           "active",
			ComplianceStatus: "non_compliant",
			LastSeen:         time.Now().UTC().Add(-2 * time.Hour),
			Metadata:         map[string]interface{}{},
		},
	}
}

func generateMockUsers() []CanonicalUser {
	return []CanonicalUser{
		{
			ID:            "user-001",
			ExternalID:    "mock-user-001",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			Email:         "john.doe@example.com",
			DisplayName:   "John Doe",
			Department:    "Engineering",
			JobTitle:      "Senior Developer",
			Status:        "active",
			Roles:         []string{"developer"},
		},
		{
			ID:            "user-002",
			ExternalID:    "mock-user-002",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			Email:         "jane.smith@example.com",
			DisplayName:   "Jane Smith",
			Department:    "Finance",
			JobTitle:      "Financial Analyst",
			Status:        "active",
			Roles:         []string{"analyst"},
		},
		{
			ID:            "user-003",
			ExternalID:    "mock-user-003",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			Email:         "admin@example.com",
			DisplayName:   "System Admin",
			Department:    "Operations",
			JobTitle:      "SysAdmin",
			Status:        "active",
			Roles:         []string{"admin"},
		},
	}
}

func init() {
	RegisterConnector("mock", func() Connector {
		return NewMockConnector()
	})
}
