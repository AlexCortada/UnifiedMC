package main

import (
	"fmt"
	"os"
	"time"

	intune "github.com/unifiedmc/connectors/intune"
	sdk "github.com/unifiedmc/connectors/sdk"
)

func main() {
	fmt.Println("=== Intune Connector Test ===")

	cfg, err := intune.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Create connector using the SDK registry
	connector, err := sdk.GlobalRegistry.Create("microsoft_intune")
	if err != nil {
		fmt.Printf("Registry error: %v\n", err)
		os.Exit(1)
	}

	// Initialize with credentials
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

	// Test connection
	fmt.Println("Testing connection to Microsoft Graph API...")

	if err := connector.Connect(); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	// Get capabilities
	caps := connector.GetCapabilities()
	fmt.Printf("Capabilities: %v\n", caps.Actions)

	// Fetch devices
	fmt.Println("Fetching devices from Intune...")
	devices, _, total, err := connector.GetDevices("", 100, nil)
	if err != nil {
		fmt.Printf("Error fetching devices: %v\n", err)
	} else {
		fmt.Printf("Found %d devices:\n", total)
		for i, d := range devices {
			if i >= 10 {
				fmt.Printf("  ... and %d more\n", total-10)
				break
			}
			fmt.Printf("  - %s (%s %s) [%s] %s\n", d.CanonicalName, d.OSType, d.OSVersion, d.ComplianceStatus, d.SerialNumber)
		}
	}

	// Fetch users
	fmt.Println("Fetching users from Entra ID...")
	users, _, userCount, err := connector.GetUsers("", 10, nil)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
	} else {
		fmt.Printf("Found %d users:\n", userCount)
		for i, u := range users {
			if i >= 10 {
				fmt.Printf("  ... and %d more\n", userCount-10)
				break
			}
			fmt.Printf("  - %s (%s) [%s]\n", u.DisplayName, u.Email, u.Department)
		}
	}

	// Health check
	health := connector.HealthCheck()
	fmt.Printf("\nHealth: %s (latency: %dms)\n", health.Status, health.LatencyMs)

	fmt.Println("\n=== Test Complete ===")
}
