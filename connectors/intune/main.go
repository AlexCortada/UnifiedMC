package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/unifiedmc/connectors/intune"
)

func main() {
	fmt.Println("=== Intune Connector Test ===")

	// Create connector
	connector := intune.NewConnector()

	// Initialize with credentials from environment
	config := sdk.ConnectorConfig{
		ConnectorType: "microsoft_intune",
		TenantID:      "test-tenant",
		Name:          "Intune Test",
		Auth: map[string]interface{}{
			"tenant_id":     os.Getenv("INTUNE_TENANT_ID"),
			"client_id":     os.Getenv("INTUNE_CLIENT_ID"),
			"client_secret": os.Getenv("INTUNE_CLIENT_SECRET"),
		},
		RateLimit: 60,
	}

	if err := connector.Initialize(config); err != nil {
		fmt.Printf("Init error (expected if no creds): %v\n", err)
		os.Exit(1)
	}

	// Test connection
	fmt.Println("Testing connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := connector.Connect(); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	// Get capabilities
	caps := connector.GetCapabilities()
	fmt.Printf("Capabilities: %v\n", caps.Actions)

	// Fetch devices
	fmt.Println("Fetching devices...")
	devices, _, total, err := connector.GetDevices("", 100, nil)
	if err != nil {
		fmt.Printf("Error fetching devices: %v\n", err)
	} else {
		fmt.Printf("Found %d devices:\n", total)
		for i, d := range devices {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", total-5)
				break
			}
			fmt.Printf("  - %s (%s %s) [%s]\n", d.CanonicalName, d.OSType, d.OSVersion, d.ComplianceStatus)
		}
	}

	// Fetch users
	fmt.Println("Fetching users...")
	users, _, userCount, err := connector.GetUsers("", 10, nil)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
	} else {
		fmt.Printf("Found %d users:\n", userCount)
		for i, u := range users {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", userCount-5)
				break
			}
			fmt.Printf("  - %s (%s) [%s]\n", u.DisplayName, u.Email, u.Department)
		}
	}

	// Health check
	health := connector.Connector.HealthCheck()
	fmt.Printf("Health: %s (latency: %dms)\n", health.Status, health.LatencyMs)

	fmt.Println("\\n=== Test Complete ===")
}
