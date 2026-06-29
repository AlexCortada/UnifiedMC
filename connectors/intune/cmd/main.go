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
			} else {
				fmt.Printf("  Error storing device %s: %v\n", d.CanonicalName, err)
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
	if metadataJSON == nil {
		metadataJSON = []byte("{}")
	}

	// Build merged sources array
	mergedSources := "{" + d.ConnectorType + "}"

	// Handle empty IP/MAC values (inet type doesn't accept empty strings)
	ipAddress := sql.NullString{String: d.IPAddress, Valid: d.IPAddress != "" && d.IPAddress != " "}
	macAddress := sql.NullString{String: d.MACAddress, Valid: d.MACAddress != "" && d.MACAddress != " "}

	// Get or create the default tenant
	tenantID := getOrCreateTenant(db, d.TenantID)

	_, err := db.Exec(`
		INSERT INTO unified_devices (
			tenant_id, display_name, asset_type, os_type, os_version,
			serial_number, manufacturer, model, mac_address, ip_address,
			status, compliance_status, last_seen,
			merged_from_sources, merge_confidence, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13,
			$14::varchar[], $15, $16::jsonb, NOW(), NOW()
		)
	`,
		tenantID,
		d.CanonicalName,
		d.AssetType,
		d.OSType,
		d.OSVersion,
		d.SerialNumber,
		d.Manufacturer,
		d.Model,
		macAddress,
		ipAddress,
		d.Status,
		d.ComplianceStatus,
		d.LastSeen,
		mergedSources,
		1.0,
		metadataJSON,
	)
	return err
}

// getOrCreateTenant returns the tenant ID, creating a default tenant if needed
func getOrCreateTenant(db *sql.DB, name string) string {
	// Try to find existing tenant
	var id string
	err := db.QueryRow("SELECT id FROM tenants WHERE name = $1", name).Scan(&id)
	if err == nil {
		return id
	}

	// Create default tenant
	if name == "" {
		name = "Default Tenant"
	}
	err = db.QueryRow(
		"INSERT INTO tenants (name, domain, status) VALUES ($1, $2, 'active') RETURNING id",
		name, name+".local",
	).Scan(&id)
	if err != nil {
		// Fallback: return first tenant
		db.QueryRow("SELECT id FROM tenants LIMIT 1").Scan(&id)
	}
	return id
}

func storeUser(db *sql.DB, u sdk.CanonicalUser) error {
	tenantID := getOrCreateTenant(db, u.TenantID)

	_, err := db.Exec(`
		INSERT INTO users (
			tenant_id, email, display_name, first_name, last_name,
			status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
		)
		ON CONFLICT (tenant_id, email) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			status = EXCLUDED.status,
			updated_at = NOW()
	`,
		tenantID,
		u.Email,
		u.DisplayName,
		u.FirstName,
		u.LastName,
		u.Status,
	)
	return err
}
