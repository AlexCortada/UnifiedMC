package intune

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// GraphClient handles authenticated requests to Microsoft Graph API
type GraphClient struct {
	cfg        *Config
	token      string
	tokenExp   time.Time
	mu         sync.RWMutex
	httpClient *http.Client
}

// NewGraphClient creates a new Graph API client
func NewGraphClient(cfg *Config) *GraphClient {
	return &GraphClient{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// authenticate obtains an OAuth2 token using client credentials flow
func (g *GraphClient) authenticate(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Return cached token if still valid
	if g.token != "" && !IsExpired(g.tokenExp) {
		return nil
	}

	data := url.Values{}
	data.Set("client_id", g.cfg.ClientID)
	data.Set("client_secret", g.cfg.ClientSecret)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", g.cfg.GetTokenURL(), strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("token request failed (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	g.token = result.AccessToken
	g.tokenExp = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return nil
}

// Get performs an authenticated GET request to the Graph API
func (g *GraphClient) Get(ctx context.Context, endpoint string) ([]byte, error) {
	if err := g.authenticate(ctx); err != nil {
		return nil, err
	}

	g.mu.RLock()
	token := g.token
	g.mu.RUnlock()

	url := g.cfg.GetGraphBaseURL() + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("ConsistencyLevel", "eventual")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

// GetPages fetches all pages of a paginated Graph API endpoint
func (g *GraphClient) GetPages(ctx context.Context, endpoint string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	nextLink := g.cfg.GetGraphBaseURL() + endpoint

	for nextLink != "" {
		if err := g.authenticate(ctx); err != nil {
			return nil, err
		}

		g.mu.RLock()
		token := g.token
		g.mu.RUnlock()

		req, err := http.NewRequestWithContext(ctx, "GET", nextLink, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("ConsistencyLevel", "eventual")

		resp, err := g.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return allItems, fmt.Errorf("graph API error (%d): %s", resp.StatusCode, string(body))
		}

		var page GraphPagedResponse
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allItems = append(allItems, page.Value...)
		nextLink = page.NextLink
	}

	return allItems, nil
}

// GetDevices retrieves all managed devices from Intune
func (g *GraphClient) GetDevices(ctx context.Context) ([]GraphDevice, error) {
	body, err := g.Get(ctx, "/deviceManagement/managedDevices?$top=999")
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []GraphDevice `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	return response.Value, nil
}

// GetUsers retrieves all users from Entra ID
func (g *GraphClient) GetUsers(ctx context.Context, top int) ([]GraphUser, error) {
	endpoint := fmt.Sprintf("/users?$top=%d&$select=id,displayName,givenName,surname,mail,userPrincipalName,jobTitle,department,accountEnabled,createdDateTime", top)
	body, err := g.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []GraphUser `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	return response.Value, nil
}

// GetApplications retrieves all managed apps from Intune
func (g *GraphClient) GetApplications(ctx context.Context) ([]GraphApp, error) {
	body, err := g.Get(ctx, "/deviceManagement/mobileApps?$top=999")
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []interface{} `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse apps: %w", err)
	}

	// Map to our simpler structure
	var apps []GraphApp
	for _, v := range response.Value {
		raw, _ := json.Marshal(v)
		var app GraphApp
		json.Unmarshal(raw, &app)
		apps = append(apps, app)
	}

	return apps, nil
}

// TestConnection verifies the connector can reach Graph API
func (g *GraphClient) TestConnection(ctx context.Context) error {
	body, err := g.Get(ctx, "/organization")
	if err != nil {
		return err
	}

	var result struct {
		DisplayName string `json:"displayName"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("organization endpoint returned unexpected data")
	}

	fmt.Printf("Connected to tenant: %s\\n", result.DisplayName)
	return nil
}
