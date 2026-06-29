package intune

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Config holds the Intune connector configuration
type Config struct {
	TenantID     string `json:"tenant_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Environment  string `json:"environment"` // azurepublic, azureusgovernment, etc.
}

// LoadConfig loads configuration from environment or JSON file
func LoadConfig() (*Config, error) {
	// Load from environment first, fall back to defaults
	cfg := &Config{
		TenantID:     getEnv("INTUNE_TENANT_ID", ""),
		ClientID:     getEnv("INTUNE_CLIENT_ID", ""),
		ClientSecret: getEnv("INTUNE_CLIENT_SECRET", ""),
		Environment:  getEnv("INTUNE_ENVIRONMENT", "azurepublic"),
	}

	if cfg.TenantID == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		// Try loading from file
		data, err := os.ReadFile("/etc/unifiedmc/intune.json")
		if err == nil {
			json.Unmarshal(data, cfg)
		}
	}

	if cfg.TenantID == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("missing Intune credentials: set INTUNE_TENANT_ID, INTUNE_CLIENT_ID, INTUNE_CLIENT_SECRET environment variables or create /etc/unifiedmc/intune.json")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func (c *Config) GetTokenURL() string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", c.TenantID)
}

func (c *Config) GetGraphBaseURL() string {
	return "https://graph.microsoft.com/v1.0"
}

// IsExpired checks if a timestamp is within the next 5 minutes
func IsExpired(t time.Time) bool {
	return time.Now().Add(5 * time.Minute).After(t)
}

// Device represents a raw Microsoft Graph managed device
type GraphDevice struct {
	ID                     string `json:"id"`
	DeviceName             string `json:"deviceName"`
	ManagedDeviceOwnerType string `json:"managedDeviceOwnerType"`
	EnrolledDateTime       string `json:"enrolledDateTime"`
	LastSyncDateTime       string `json:"lastSyncDateTime"`
	OperatingSystem        string `json:"operatingSystem"`
	OSVersion              string `json:"osVersion"`
	ComplianceState        string `json:"complianceState"`
	JailBroken             string `json:"jailBroken"`
	ManagementAgent        string `json:"managementAgent"`
	EnrollmentProfileName  string `json:"enrollmentProfileName"`
	UserPrincipalName      string `json:"userPrincipalName"`
	UserID                 string `json:"userId"`
	Model                  string `json:"model"`
	Manufacturer           string `json:"manufacturer"`
	SerialNumber           string `json:"serialNumber"`
	PhoneNumber            string `json:"phoneNumber"`
	WiFiMacAddress         string `json:"wifiMacAddress"`
	IMEI                   string `json:"imei"`
	StorageTotal           int64  `json:"totalStorageSpaceInBytes"`
	StorageFree            int64  `json:"freeStorageSpaceInBytes"`
	IsEncrypted            bool   `json:"isEncrypted"`
	IsSupervised           bool   `json:"isSupervised"`
	AADRegistered          bool   `json:"isAzureADRegistered"`
	ManagementCertificates string `json:"managementCertificateExpiryDate"`
	ETag                   string `json:"eTag"`
	DeviceCategory         struct {
		DisplayName string `json:"displayName"`
	} `json:"deviceCategory"`
}

// User represents a raw Microsoft Graph user
type GraphUser struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	JobTitle          string `json:"jobTitle"`
	Department        string `json:"department"`
	Manager           struct {
		ID string `json:"id"`
	} `json:"manager"`
	AccountEnabled    bool     `json:"accountEnabled"`
	CreatedDateTime   string   `json:"createdDateTime"`
	LastSignInDateTime string  `json:"lastSignInDateTime"`
	MemberOf          []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"memberOf"`
}

// Application represents a raw Microsoft Graph mobile app
type GraphApp struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayVersion"`
	Description     string `json:"description"`
	Publisher       string `json:"publisher"`
	Version         string `json:"version"`
	BundleID        string `json:"bundleId"`
	ExpirationDateTime string `json:"expirationDateTime"`
}

// GraphPagedResponse holds paginated Graph API responses
type GraphPagedResponse struct {
	Value    json.RawMessage `json:"value"`
	NextLink string          `json:"@odata.nextLink"`
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}
