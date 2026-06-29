package main

import (
	"fmt"
	"os"

	"github.com/unifiedmc/connectors/intune"
	_ "github.com/unifiedmc/connectors/intune"
	sdk "github.com/unifiedmc/connectors/sdk"
)

func main() {
	fmt.Println("=== Intune Connector Test ===")

	// Use LoadConfig which checks env vars then falls back to /etc/unifiedmc/intune.json
	cfg, err := intune.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		fmt.Println("\nTroubleshooting:")
		fmt.Println("  1. Set environment variables: INTUNE_TENANT_ID, INTUNE_CLIENT_ID, INTUNE_CLIENT_SECRET")
		fmt.Println("  2. Or create /etc/unifiedmc/intune.json with tenant_id, client_id, client_secret")
		os.Exit(1)
	}

	connector, err := sdk.GlobalRegistry.Create("microsoft_intune")
	if err != nil {
		fmt.Printf("Registry error: %v\n", err)
		os.Exit(1)
	}

	config := sdk.ConnectorConfig{
		ConnectorType: "microsoft_intune",
		TenantID:      cfg.TenantID,
		Name:          "Intune Production",
		Auth: map[string]interface{}{
			"tenant_id":     cfg.TenantID,
			"client_id":     cfg.ClientID,
			"client_secret": cfg.ClientSecret,
		},
		RateLimit: 60,
	}

	if err := connector.Initialize(config); err != nil {
		fmt.Printf("Init error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Tenant: %s\n", cfg.TenantID)
	fmt.Printf("Client: %s\n", cfg.ClientID)
	fmt.Printf("Secret length: %d\n", len(cfg.ClientSecret))

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
