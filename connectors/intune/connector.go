package intune

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/unifiedmc/connectors/sdk"
)

// Connector implements the IConnector interface for Microsoft Intune
type Connector struct {
	sdk.BaseConnector
	client *GraphClient
}

// NewConnector creates a new Intune connector
func NewConnector() *Connector {
	return &Connector{}
}

// Initialize sets up the connector with configuration
func (c *Connector) Initialize(config sdk.ConnectorConfig) error {
	c.BaseInitialize(config, "intune")

	// Build Intune-specific config from connector config
	intuneCfg := &Config{
		TenantID:     getString(config.Auth, "tenant_id"),
		ClientID:     getString(config.Auth, "client_id"),
		ClientSecret: getString(config.Auth, "client_secret"),
		Environment:  getString(config.Metadata, "environment"),
	}

	c.client = NewGraphClient(intuneCfg)
	return nil
}

// Connect authenticates and verifies connectivity
func (c *Connector) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.client.TestConnection(ctx); err != nil {
		return fmt.Errorf("intune connection failed: %w", err)
	}

	c.BaseConnect()
	return nil
}

// Disconnect cleans up resources
func (c *Connector) Disconnect() error {
	c.BaseDisconnect()
	return nil
}

// HealthCheck returns the connector's health status
func (c *Connector) HealthCheck() sdk.HealthStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	err := c.client.TestConnection(ctx)

	return sdk.HealthStatus{
		ConnectorType: "microsoft_intune",
		Status:        connectionStatus(err == nil),
		LatencyMs:     int(time.Since(start).Milliseconds()),
		ErrorMessage:  errString(err),
	}
}

// GetCapabilities declares what this connector supports
func (c *Connector) GetCapabilities() sdk.ConnectorCapabilities {
	return sdk.ConnectorCapabilities{
		ConnectorType: "microsoft_intune",
		EntityTypes:   []string{"device", "user", "application"},
		Operations:    []string{"read", "sync"},
		Actions:       []string{"run_script", "restart_device", "deploy_application", "patch_device"},
		AuthMethods:   []string{"oauth2"},
		SyncModes:     []string{"pull"},
		RateLimit:     60,
	}
}

// GetDevices retrieves all managed devices from Intune
func (c *Connector) GetDevices(cursor string, pageSize int, filters map[string]interface{}) ([]sdk.CanonicalDevice, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	graphDevices, err := c.client.GetDevices(ctx)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to fetch devices: %w", err)
	}

	devices := make([]sdk.CanonicalDevice, 0, len(graphDevices))
	for _, d := range graphDevices {
		devices = append(devices, mapToCanonicalDevice(d, c.Config.TenantID))
	}

	return devices, "", len(devices), nil
}

// GetDevice retrieves a single device by external ID
func (c *Connector) GetDevice(deviceID string) (*sdk.CanonicalDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch all and filter (Graph API doesn't support direct ID lookup for managedDevices easily)
	devices, err := c.client.GetDevices(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.ID == deviceID {
			dev := mapToCanonicalDevice(d, c.Config.TenantID)
			return &dev, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", deviceID)
}

// GetUsers retrieves all users from Entra ID
func (c *Connector) GetUsers(cursor string, pageSize int, filters map[string]interface{}) ([]sdk.CanonicalUser, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	top := 999
	if pageSize > 0 {
		top = pageSize
	}

	graphUsers, err := c.client.GetUsers(ctx, top)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to fetch users: %w", err)
	}

	users := make([]sdk.CanonicalUser, 0, len(graphUsers))
	for _, u := range graphUsers {
		users = append(users, mapToCanonicalUser(u, c.Config.TenantID))
	}

	return users, "", len(users), nil
}

// RunScript executes a script on a managed device
func (c *Connector) RunScript(deviceID string, content string, scriptType string, timeout int) (sdk.ScriptResult, error) {
	return sdk.ScriptResult{
		DeviceID: deviceID,
		Status:   "not_implemented",
		Stdout:   "Script execution not yet implemented for Intune",
	}, nil
}

// RestartDevice reboots a managed device
func (c *Connector) RestartDevice(deviceID string, force bool, reason string) (sdk.ActionResult, error) {
	return sdk.ActionResult{
		DeviceID:   deviceID,
		ActionType: "restart",
		Status:     "not_implemented",
		Message:    "Restart not yet implemented for Intune",
	}, nil
}

// DeployApplication deploys an app to devices
func (c *Connector) DeployApplication(appID string, deviceIDs []string) (sdk.DeploymentResult, error) {
	return sdk.DeploymentResult{
		DeviceID: deviceIDs[0],
		Status:   "not_implemented",
		Message:  "Deployment not yet implemented for Intune",
	}, nil
}

// --- Mapping Functions ---

func mapToCanonicalDevice(d GraphDevice, tenantID string) sdk.CanonicalDevice {
	device := sdk.CanonicalDevice{
		ExternalID:    d.ID,
		ConnectorType: "microsoft_intune",
		TenantID:      tenantID,
		CanonicalName: d.DeviceName,
		OSType:        mapOSType(d.OperatingSystem),
		OSVersion:     d.OSVersion,
		SerialNumber:  d.SerialNumber,
		Manufacturer:  d.Manufacturer,
		Model:         d.Model,
		MACAddress:    d.WiFiMacAddress,
		PrimaryUserID: d.UserID,
		Status:        mapDeviceStatus(d),
		Metadata: map[string]interface{}{
			"intune_enrolled":      d.EnrolledDateTime,
			"intune_last_sync":     d.LastSyncDateTime,
			"intune_compliance":    d.ComplianceState,
			"intune_is_encrypted":  d.IsEncrypted,
			"intune_is_supervised": d.IsSupervised,
			"intune_agent":         d.ManagementAgent,
			"intune_owner_type":    d.ManagedDeviceOwnerType,
		},
	}

	// Map compliance status
	switch strings.ToLower(d.ComplianceState) {
	case "compliant":
		device.ComplianceStatus = "compliant"
	case "noncompliant", "conflict", "error":
		device.ComplianceStatus = "non_compliant"
	default:
		device.ComplianceStatus = "unknown"
	}

	// Map asset type
	device.AssetType = mapAssetType(d.OperatingSystem, d.Model)

	// Parse last seen
	if d.LastSyncDateTime != "" {
		if t, err := time.Parse(time.RFC3339, d.LastSyncDateTime); err == nil {
			device.LastSeen = t
		}
	}

	return device
}

func mapToCanonicalUser(u GraphUser, tenantID string) sdk.CanonicalUser {
	user := sdk.CanonicalUser{
		ExternalID:    u.ID,
		ConnectorType: "microsoft_entra_id",
		TenantID:      tenantID,
		Email:         u.Mail,
		DisplayName:   u.DisplayName,
		FirstName:     u.GivenName,
		LastName:      u.Surname,
		Department:    u.Department,
		JobTitle:      u.JobTitle,
		Status:        mapUserStatus(u.AccountEnabled),
	}

	if u.Manager.ID != "" {
		user.ManagerID = u.Manager.ID
	}

	return user
}

func mapOSType(os string) string {
	switch strings.ToLower(os) {
	case "windows":
		return "windows"
	case "macos", "mac", "mac_md":
		return "macos"
	case "ios":
		return "ios"
	case "android":
		return "android"
	case "linux":
		return "linux"
	default:
		return "other"
	}
}

func mapAssetType(os string, model string) string {
	osLower := strings.ToLower(os)
	if osLower == "ios" || osLower == "android" {
		return "mobile"
	}
	if strings.Contains(strings.ToLower(model), "virtual") {
		return "vm"
	}
	return "workstation"
}

func mapDeviceStatus(d GraphDevice) string {
	if d.LastSyncDateTime == "" {
		return "inactive"
	}
	// If last sync was within 7 days, consider active
	if t, err := time.Parse(time.RFC3339, d.LastSyncDateTime); err == nil {
		if time.Since(t) < 7*24*time.Hour {
			return "active"
		}
	}
	return "inactive"
}

func mapUserStatus(enabled bool) string {
	if enabled {
		return "active"
	}
	return "disabled"
}

func connectionStatus(healthy bool) string {
	if healthy {
		return "connected"
	}
	return "disconnected"
}

func errString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func init() {
	sdk.RegisterConnector("microsoft_intune", func() sdk.Connector {
		return NewConnector()
	})
}
