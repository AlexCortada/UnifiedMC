package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// CanonicalDevice represents a unified device
type CanonicalDevice struct {
	ID               string `json:"id"`
	ExternalID       string `json:"external_id"`
	ConnectorType    string `json:"connector_type"`
	TenantID         string `json:"tenant_id"`
	CanonicalName    string `json:"canonical_name"`
	AssetType        string `json:"asset_type"`
	OSType           string `json:"os_type"`
	OSVersion        string `json:"os_version"`
	SerialNumber     string `json:"serial_number"`
	ComplianceStatus string `json:"compliance_status"`
	Status           string `json:"status"`
	LastSeen         string `json:"last_seen"`
}

// DeviceResponse is the API response for device endpoints
type DeviceResponse struct {
	Devices []CanonicalDevice `json:"devices"`
	Total   int               `json:"total"`
	Source  string            `json:"source"`
}

// HealthResponse is the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Version   string `json:"version"`
}

// GetEnv retrieves environment variable with fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HealthHandler returns health status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   GetEnv("SERVICE_NAME", "uop-api-gateway"),
		Version:   GetEnv("VERSION", "0.1.0"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DevicesHandler returns mock device list
func DevicesHandler(w http.ResponseWriter, r *http.Request) {
	// In production, this would call the inventory service
	devices := []CanonicalDevice{
		{
			ID:            "dev-001",
			ExternalID:    "mock-001",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			CanonicalName: "DESKTOP-ABC123",
			AssetType:     "workstation",
			OSType:        "windows",
			Status:        "active",
			LastSeen:      time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339),
		},
		{
			ID:            "dev-002",
			ExternalID:    "mock-002",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			CanonicalName: "LAPTOP-XYZ789",
			AssetType:     "workstation",
			OSType:        "macos",
			Status:        "active",
			LastSeen:      time.Now().UTC().Add(-12 * time.Minute).Format(time.RFC3339),
		},
		{
			ID:            "dev-003",
			ExternalID:    "mock-003",
			ConnectorType: "mock",
			TenantID:      "tenant-001",
			CanonicalName: "SERVER-WEB01",
			AssetType:     "server",
			OSType:        "linux",
			Status:        "active",
			LastSeen:      time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339),
		},
	}

	response := DeviceResponse{
		Devices: devices,
		Total:   len(devices),
		Source:  "mock",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeviceDetailHandler returns a single device
func DeviceDetailHandler(w http.ResponseWriter, r *http.Request) {
	device := CanonicalDevice{
		ID:               "dev-001",
		ExternalID:       "mock-001",
		ConnectorType:    "mock",
		TenantID:         "tenant-001",
		CanonicalName:    "DESKTOP-ABC123",
		AssetType:        "workstation",
		OSType:           "windows",
		OSVersion:        "11 23H2",
		SerialNumber:     "SN-ABC-12345",
		ComplianceStatus: "compliant",
		Status:           "active",
		LastSeen:         time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	log.Println("Unified IT Operations Portal - API Gateway Starting...")
	log.Printf("Service: %s", GetEnv("SERVICE_NAME", "uop-api-gateway"))
	log.Printf("Port: %s", GetEnv("PORT", "8080"))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/v1/devices", DevicesHandler)
	mux.HandleFunc("/api/v1/devices/", DeviceDetailHandler)

	port := GetEnv("PORT", "8080")
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, LoggingMiddleware(mux)))
}
