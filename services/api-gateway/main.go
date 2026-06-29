package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// CanonicalDevice represents a unified device
type CanonicalDevice struct {
	ID               string    `json:"id"`
	ExternalID       string    `json:"external_id,omitempty"`
	ConnectorType    string    `json:"connector_type"`
	TenantID         string    `json:"tenant_id"`
	CanonicalName    string    `json:"canonical_name"`
	AssetType        string    `json:"asset_type"`
	OSType           string    `json:"os_type"`
	OSVersion        string    `json:"os_version,omitempty"`
	SerialNumber     string    `json:"serial_number,omitempty"`
	Manufacturer     string    `json:"manufacturer,omitempty"`
	Model            string    `json:"model,omitempty"`
	Status           string    `json:"status"`
	ComplianceStatus string    `json:"compliance_status"`
	LastSeen         string    `json:"last_seen,omitempty"`
	PatchStatus      string    `json:"patch_status,omitempty"`
}

// DeviceResponse is the API response
type DeviceResponse struct {
	Devices []CanonicalDevice `json:"devices"`
	Total   int               `json:"total"`
	Source  string            `json:"source"`
}

// HealthResponse is the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"version"`
	Version   string `json:"version"`
}

var db *sql.DB

func main() {
	log.Println("Unified IT Operations Portal - API Gateway Starting...")

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://unifiedmc:***@127.0.0.1:5432/unifiedmc?sslmode=disable"
	}
	dbURL = strings.TrimSpace(dbURL)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/devices", devicesHandler)
	mux.HandleFunc("/api/v1/devices/", deviceDetailHandler)
	mux.HandleFunc("/api/v1/dashboard", dashboardHandler)
	mux.HandleFunc("/api/v1/sync/trigger", syncTriggerHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, loggingMiddleware(mux)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "uop-api-gateway",
		Version:   "0.2.0",
	}
	if db != nil {
		if err := db.Ping(); err != nil {
			response.Status = "degraded"
		}
	}
	writeJSON(w, response)
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := getDevicesFromDB(r.URL.Query().Get("os_type"), r.URL.Query().Get("status"), r.URL.Query().Get("compliance"))
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		http.Error(w, `{"error": "database error"}`, 500)
		return
	}

	writeJSON(w, DeviceResponse{
		Devices: devices,
		Total:   getDeviceCount(),
		Source:  "database",
	})
}

func getDeviceCount() int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM unified_devices").Scan(&count)
	return count
}

func deviceDetailHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := strings.TrimPrefix(r.URL.Path, "/api/v1/devices/")
	if deviceID == "" {
		http.Error(w, `{"error": "device ID required"}`, 400)
		return
	}

	var device CanonicalDevice
	var displayName sql.NullString
	var lastSeen sql.NullTime
	err := db.QueryRow(`
		SELECT id, display_name, asset_type, os_type, os_version, serial_number,
		       manufacturer, model, status, compliance_status, last_seen
		FROM unified_devices WHERE id = $1 OR display_name ILIKE $2
		LIMIT 1`, deviceID, "%"+deviceID+"%").Scan(
		&device.ID, &displayName, &device.AssetType, &device.OSType,
		&device.OSVersion, &device.SerialNumber, &device.Manufacturer, &device.Model,
		&device.Status, &device.ComplianceStatus, &device.LastSeen,
	)
	if err != nil {
		http.Error(w, `{"error": "device not found"}`, 404)
		return
	}
	if displayName.Valid {
		device.CanonicalName = displayName.String
	}
	if lastSeen.Valid {
		device.LastSeen = lastSeen.Time.Format(time.RFC3339)
	}

	writeJSON(w, device)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	summary := getDashboardSummary()
	writeJSON(w, summary)
}

func syncTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "method not allowed"}`, 405)
		return
	}
	writeJSON(w, map[string]string{"status": "sync_queued", "message": "Sync triggered via API"})
}

func getDevicesFromDB(osType, status, compliance string) ([]CanonicalDevice, error) {
	query := `
		SELECT id, display_name, asset_type, os_type, os_version, serial_number,
		       manufacturer, model, status, compliance_status, last_seen
		FROM unified_devices
		WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if osType != "" {
		query += fmt.Sprintf(" AND os_type = $%d", argNum)
		args = append(args, osType)
		argNum++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
		argNum++
	}
	if compliance != "" {
		query += fmt.Sprintf(" AND compliance_status = $%d", argNum)
		args = append(args, compliance)
		argNum++
	}

	query += " ORDER BY display_name LIMIT 10000"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []CanonicalDevice
	for rows.Next() {
		var d CanonicalDevice
		var displayName sql.NullString
		var lastSeen sql.NullTime
		err := rows.Scan(
			&d.ID, &displayName, &d.AssetType, &d.OSType,
			&d.OSVersion, &d.SerialNumber, &d.Manufacturer, &d.Model,
			&d.Status, &d.ComplianceStatus, &lastSeen,
		)
		if err != nil {
			return nil, err
		}
		if displayName.Valid {
			d.CanonicalName = displayName.String
		}
		if lastSeen.Valid {
			d.LastSeen = lastSeen.Time.Format(time.RFC3339)
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func getDashboardSummary() map[string]interface{} {
	var total, online, compliant, nonCompliant int

	db.QueryRow("SELECT COUNT(*) FROM unified_devices").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM unified_devices WHERE status = 'active'").Scan(&online)
	db.QueryRow("SELECT COUNT(*) FROM unified_devices WHERE compliance_status = 'compliant'").Scan(&compliant)
	db.QueryRow("SELECT COUNT(*) FROM unified_devices WHERE compliance_status = 'non_compliant'").Scan(&nonCompliant)

	rows, _ := db.Query("SELECT os_type, COUNT(*) FROM unified_devices GROUP BY os_type ORDER BY os_type")
	osBreakdown := map[string]int{}
	for rows.Next() {
		var os string
		var count int
		rows.Scan(&os, &count)
		osBreakdown[os] = count
	}
	rows.Close()

	return map[string]interface{}{
		"total_devices":      total,
		"online_devices":     online,
		"offline_devices":    total - online,
		"compliant_count":    compliant,
		"non_compliant_count": nonCompliant,
		"compliance_rate":    fmt.Sprintf("%.1f", float64(compliant)/float64(total)*100),
		"os_breakdown":       osBreakdown,
		"source":             "database",
		"last_updated":       time.Now().UTC().Format(time.RFC3339),
	}
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
