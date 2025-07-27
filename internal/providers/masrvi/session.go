package masrvi

import (
	"context"
	"fmt"
	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SessionManager handles MASRVI session management
type SessionManager struct {
	config     rimpay.ProviderConfig
	httpClient common.HTTPClient
	logger     rimpay.Logger
	baseURL    string

	// Session cache
	sessionCache map[string]*sessionCacheEntry
	cacheMutex   sync.RWMutex
}

type sessionCacheEntry struct {
	sessionID string
	expiresAt time.Time
}

// NewSessionManager creates new session manager
func NewSessionManager(config rimpay.ProviderConfig, httpClient common.HTTPClient, logger rimpay.Logger) *SessionManager {
	return &SessionManager{
		config:       config,
		httpClient:   httpClient,
		logger:       logger,
		baseURL:      strings.TrimRight(config.BaseURL, "/"),
		sessionCache: make(map[string]*sessionCacheEntry),
	}
}

// GetSessionID gets a valid session ID
func (sm *SessionManager) GetSessionID(ctx context.Context) (string, error) {
	merchantID := sm.config.Credentials["merchant_id"]

	// Check cache first
	sm.cacheMutex.RLock()
	if entry, exists := sm.sessionCache[merchantID]; exists && time.Now().Before(entry.expiresAt) {
		sessionID := entry.sessionID
		sm.cacheMutex.RUnlock()
		sm.logger.Debug("Using cached session ID", "session_id", sessionID)
		return sessionID, nil
	}
	sm.cacheMutex.RUnlock()

	// Get new session
	return sm.createSession(ctx, merchantID)
}

// createSession creates a new session
func (sm *SessionManager) createSession(ctx context.Context, merchantID string) (string, error) {
	sessionURL := fmt.Sprintf("%s/online/online.php?merchantid=%s", sm.baseURL, merchantID)

	req := &common.HTTPRequest{
		Method:  "GET",
		URL:     sessionURL,
		Headers: make(map[string]string),
		Timeout: sm.config.Timeout,
	}

	sm.logger.Debug("Creating MASRVI session", "merchant_id", merchantID)

	resp, err := sm.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("session creation failed with status: %d", resp.StatusCode)
	}

	sessionID := strings.TrimSpace(string(resp.Body))
	if sessionID == "" || sessionID == "NOK" {
		return "", fmt.Errorf("invalid session response: %s", sessionID)
	}

	// Cache the session (TTL: 5 minutes)
	sm.cacheMutex.Lock()
	sm.sessionCache[merchantID] = &sessionCacheEntry{
		sessionID: sessionID,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	sm.cacheMutex.Unlock()

	sm.logger.Info("MASRVI session created", "session_id", sessionID)

	return sessionID, nil
}

// ClearCache clears the session cache
func (sm *SessionManager) ClearCache() {
	sm.cacheMutex.Lock()
	defer sm.cacheMutex.Unlock()
	sm.sessionCache = make(map[string]*sessionCacheEntry)
}
