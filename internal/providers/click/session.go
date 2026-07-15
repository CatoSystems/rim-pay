package click

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// SessionManager handles CLICK (TagPay) session IDs.
type SessionManager struct {
	config     rimpay.ProviderConfig
	httpClient common.HTTPClient
	logger     rimpay.Logger
	baseURL    string

	sessionCache map[string]*sessionCacheEntry
	cacheMutex   sync.RWMutex
}

type sessionCacheEntry struct {
	sessionID string
	expiresAt time.Time
}

// NewSessionManager creates a new CLICK session manager.
func NewSessionManager(config rimpay.ProviderConfig, httpClient common.HTTPClient, logger rimpay.Logger) *SessionManager {
	return &SessionManager{
		config:       config,
		httpClient:   httpClient,
		logger:       logger,
		baseURL:      strings.TrimRight(config.BaseURL, "/"),
		sessionCache: make(map[string]*sessionCacheEntry),
	}
}

// GetSessionID returns a valid (cached or fresh) session ID.
func (sm *SessionManager) GetSessionID(ctx context.Context) (string, error) {
	merchantID := sm.config.Credentials["merchant_id"]

	sm.cacheMutex.RLock()
	if entry, ok := sm.sessionCache[merchantID]; ok && time.Now().Before(entry.expiresAt) {
		id := entry.sessionID
		sm.cacheMutex.RUnlock()
		return id, nil
	}
	sm.cacheMutex.RUnlock()

	return sm.createSession(ctx, merchantID)
}

func (sm *SessionManager) createSession(ctx context.Context, merchantID string) (string, error) {
	sessionURL := fmt.Sprintf("%s/online/online.php?merchantid=%s", sm.baseURL, merchantID)

	resp, err := sm.httpClient.Do(&common.HTTPRequest{
		Method:  "GET",
		URL:     sessionURL,
		Headers: make(map[string]string),
		Timeout: sm.config.Timeout,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("session creation failed with status: %d", resp.StatusCode)
	}

	raw := strings.TrimSpace(string(resp.Body))
	// TagPay returns "OK:<sessionid>" or "NOK:<REASON>".
	switch {
	case strings.HasPrefix(raw, "OK:"):
		sessionID := strings.TrimPrefix(raw, "OK:")
		if sessionID == "" {
			return "", fmt.Errorf("empty session id in response: %q", raw)
		}
		sm.cacheMutex.Lock()
		sm.sessionCache[merchantID] = &sessionCacheEntry{
			sessionID: sessionID,
			expiresAt: time.Now().Add(180 * time.Second), // spec default session timeout
		}
		sm.cacheMutex.Unlock()
		sm.logger.Info("CLICK session created", "merchant_id", merchantID)
		return sessionID, nil
	case strings.HasPrefix(raw, "NOK:"):
		return "", fmt.Errorf("session refused: %s", strings.TrimPrefix(raw, "NOK:"))
	default:
		return "", fmt.Errorf("unexpected session response: %q", raw)
	}
}
