package main

import (
	"fmt"
	"os"

	_ "github.com/unifiedmc/connectors/intune"
	sdk "github.com/unifiedmc/connectors/sdk"
)

func main() {
	fmt.Println("=== Intune Connector Test ===")

	// Read from environment or use defaults for testing
	tenantID := os.Getenv("INTUNE_TENANT_ID")
	clientID := os.Getenv("INTUNE_CLIENT_ID")
	clientSecret := os.Getenv("INTUNE_CLIENT_SECRET")

	if tenantID == "" || clientID == "" || clientSecret == "" {
		fmt.Println("Missing credentials. Set environment variables:")
		fmt.Println("  INTUNE_TENANT_ID")
		fmt.Println("  INTUNE_CLIENT_ID")
		fmt.Println("  INTUNE_CLIENT_SECRET")
		os.Exit(1)
	}

	connector, err := sdk.GlobalRegistry.Create("microsoft_intune")
	if err != nil {
		fmt.Printf("Registry error: %v\n", err)
		os.Exit(1)
	}

	config := sdk.ConnectorConfig{
		ConnectorType: "microsoft_intune",
		TenantID:      tenantID,
		Name:          "Intune Production",
		Auth: map[string]interface{}{
			"tenant_id":     tenantID,
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
		RateLimit: 60,
	}

	if err := connector.Initialize(config); err != nil {
		fmt.Printf("Init error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Tenant: %s\n", tenantID)
	fmt.Printf("Client: %s\n", clientID)
	fmt.Printf("Secret length: %d\n", len(clientSecret))

	fmt.Println("\nTesting connection...")
	if err := connector.Connect(); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	devices, _, total, err := connector.GetDevices("", 100, nil)
	if err != nil {
		fmt.Printf("Error fetching devices: %v\n", err)
	} else {
		fmt.Printf("Found %d devices\n", total)
		for _, d := range devices {
			fmt.Printf("  - %s (%s) [%s]\n", d.CanonicalName, d.OSType, d.ComplianceStatus)
		}
	}
}
