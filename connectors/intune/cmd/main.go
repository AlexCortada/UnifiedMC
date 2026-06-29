package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/unifiedmc/connectors/intune"
	_ "github.com/unifiedmc/connectors/intune"
	sdk "github.com/unifiedmc/connectors/sdk"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("=== Intune Connector - Full Sync ===")

	// Load config
	cfg, err := intune.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Create connector
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

	// Test connection
	fmt.Println("Connecting to Microsoft Graph API...")
	if err := connector.Connect(); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	// Connect to database - always use default (ignore env to avoid SSH corruption)
	dbURL := "postgresql://unifiedmc:***@127.0.0.1:5432/unifiedmc?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Database unreachable: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to PostgreSQL")

	// Fetch and store devices
	fmt.Println("\nFetching devices from Intune...")
	devices, _, total, err := connector.GetDevices("", 999, nil)
	if err != nil {
		fmt.Printf("Error fetching devices: %v\n", err)
	} else {
		fmt.Printf("Found %d devices. Storing in database...\n", total)
		stored := 0
		for _, d := range devices {
			if err := storeDevice(db, d); err == nil {
				stored++
			}
		}
		fmt.Printf("Stored/updated %d devices\n", stored)
	}

	// Fetch and store users
	fmt.Println("\nFetching users from Entra ID...")
	users, _, userCount, err := connector.GetUsers("", 999, nil)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
	} else {
		fmt.Printf("Found %d users. Storing in database...\n", userCount)
		stored := 0
		for _, u := range users {
			if err := storeUser(db, u); err == nil {
				stored++
			}
		}
		fmt.Printf("Stored/updated %d users\n", stored)
	}

	// Summary
	var deviceCount, userCountDB int
	db.QueryRow("SELECT COUNT(*) FROM unified_devices").Scan(&deviceCount)
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCountDB)

	fmt.Printf("\n=== Sync Complete ===\n")
	fmt.Printf("Devices in database: %d\n", deviceCount)
	fmt.Printf("Users in database: %d\n", userCountDB)
}

func storeDevice(db *sql.DB, d sdk.CanonicalDevice) error {
	metadataJSON, _ := json.Marshal(d.Metadata)
	var lastSeen interface{}
	if !d.LastSeen.IsZero() {
		lastSeen = d.LastSeen
	}

	_, err := db.Exec(`
		INSERT INTO unified_devices (
			tenant_id, display_name, asset_type, os_type, os_version,
			serial_number, manufacturer, model, primary_user_id, ip_address,
			mac_address, status, compliance_status, last_seen, merged_from_sources,
			merge_confidence, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, NOW(), NOW()
		)
		ON CONFLICT (serial_number) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			os_version = EXCLUDED.os_version,
			status = EXCLUDED.status,
			compliance_status = EXCLUDED.compliance_status,
			last_seen = EXCLUDED.last_seen,
			updated_at = NOW()
	`,
		d.TenantID, d.CanonicalName, d.AssetType, d.OSType, d.OSVersion,
		d.SerialNumber, d.Manufacturer, d.Model, d.PrimaryUserID, d.IPAddress,
		d.MACAddress, d.Status, d.ComplianceStatus, lastSeen,
		fmt.Sprintf("{%s}", d.ConnectorType), 1.0, metadataJSON,
	)
	return err
}

func storeUser(db *sql.DB, u sdk.CanonicalUser) error {
	_, err := db.Exec(`
		INSERT INTO users (
			tenant_id, email, display_name, first_name, last_name,
			department, job_title, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		ON CONFLICT (tenant_id, email) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			department = EXCLUDED.department,
			job_title = EXCLUDED.job_title,
			status = EXCLUDED.status,
			updated_at = NOW()
	`,
		u.TenantID, u.Email, u.DisplayName, u.FirstName, u.LastName,
		u.Department, u.JobTitle, u.Status,
	)
	return err
}
