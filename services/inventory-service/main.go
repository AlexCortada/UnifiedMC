package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Device represents a unified device in the database
type Device struct {
	ID               string
	ExternalID       string
	ConnectorType    string
	TenantID         string
	CanonicalName    string
	AssetType        string
	OSType           string
	OSVersion        string
	SerialNumber     string
	Manufacturer     string
	Model            string
	PrimaryUserID    string
	IPAddress        string
	MACAddress       string
	Status           string
	ComplianceStatus string
	LastSeen         time.Time
	Metadata         map[string]string
}

// User represents a unified user in the database
type User struct {
	ID            string
	ExternalID    string
	ConnectorType string
	TenantID      string
	Email         string
	DisplayName   string
	FirstName     string
	LastName      string
	Department    string
	JobTitle      string
	ManagerID     string
	Status        string
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://unifiedmc:***@localhost:5432/unifiedmc?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Test: insert a sample device
	device := Device{
		ExternalID:       "test-001",
		ConnectorType:    "microsoft_intune",
		TenantID:         "tenant-001",
		CanonicalName:    "TEST-DEVICE",
		AssetType:        "workstation",
		OSType:           "windows",
		OSVersion:        "11 23H2",
		SerialNumber:     "TEST-SN-123",
		Manufacturer:     "Dell",
		Model:            "OptiPlex 7090",
		Status:           "active",
		ComplianceStatus: "compliant",
		LastSeen:         time.Now().UTC(),
		Metadata:         map[string]string{"source": "test"},
	}

	if err := upsertDevice(db, device); err != nil {
		log.Printf("Failed to insert device: %v", err)
	} else {
		fmt.Println("Device inserted successfully")
	}

	// Query devices
	devices, err := getDevices(db)
	if err != nil {
		log.Printf("Failed to query devices: %v", err)
	} else {
		fmt.Printf("Found %d devices in database:\n", len(devices))
		for _, d := range devices {
			fmt.Printf("  - %s (%s %s) [%s]\n", d.CanonicalName, d.OSType, d.OSVersion, d.ComplianceStatus)
		}
	}
}

func upsertDevice(db *sql.DB, d Device) error {
	metadataJSON, _ := json.Marshal(d.Metadata)

	query := `
		INSERT INTO unified_devices (
			id, tenant_id, display_name, asset_type, os_type, os_version,
			serial_number, manufacturer, model, primary_user_id, ip_address,
			mac_address, status, compliance_status, last_seen, merged_from_sources,
			merge_confidence, metadata, created_at, updated_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, NOW(), NOW()
		)
		ON CONFLICT (serial_number) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			os_version = EXCLUDED.os_version,
			status = EXCLUDED.status,
			compliance_status = EXCLUDED.compliance_status,
			last_seen = EXCLUDED.last_seen,
			updated_at = NOW()
		RETURNING id
	`

	var id string
	err := db.QueryRow(query,
		d.TenantID, d.CanonicalName, d.AssetType, d.OSType, d.OSVersion,
		d.SerialNumber, d.Manufacturer, d.Model, d.PrimaryUserID, d.IPAddress,
		d.MACAddress, d.Status, d.ComplianceStatus, d.LastSeen,
		"{" + d.ConnectorType + "}", 1.0, metadataJSON,
	).Scan(&id)

	if err == nil {
		fmt.Printf("  Device ID: %s\n", id)
	}
	return err
}

func getDevices(db *sql.DB) ([]Device, error) {
	rows, err := db.Query(`
		SELECT id, external_id, connector_type, tenant_id, display_name, 
		       asset_type, os_type, os_version, serial_number, manufacturer,
		       model, status, compliance_status, last_seen
		FROM unified_devices 
		ORDER BY display_name
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		var lastSeen sql.NullTime
		err := rows.Scan(
			&d.ID, &d.ExternalID, &d.ConnectorType, &d.TenantID, &d.CanonicalName,
			&d.AssetType, &d.OSType, &d.OSVersion, &d.SerialNumber, &d.Manufacturer,
			&d.Model, &d.Status, &d.ComplianceStatus, &lastSeen,
		)
		if err != nil {
			return nil, err
		}
		if lastSeen.Valid {
			d.LastSeen = lastSeen.Time
		}
		devices = append(devices, d)
	}
	return devices, nil
}
