package bpay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type AuthManager struct {
	config     rimpay.ProviderConfig
	httpClient common.HTTPClient
	logger     rimpay.Logger

	// Authentication state
	auth      *AuthResponse
	authMutex sync.RWMutex
	baseURL   string
}

// NewAuthManager creates new authentication manager
func NewAuthManager(config rimpay.ProviderConfig, httpClient common.HTTPClient, logger rimpay.Logger) *AuthManager {
	return &AuthManager{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
		baseURL:    strings.TrimRight(config.BaseURL, "/"),
	}
}

// GetAccessToken gets valid access token
func (am *AuthManager) GetAccessToken(ctx context.Context) (string, error) {
	am.authMutex.RLock()
	if am.auth != nil && !am.isTokenExpired() {
		token := am.auth.AccessToken
		am.authMutex.RUnlock()
		return token, nil
	}
	am.authMutex.RUnlock()

	// Token expired or not available, authenticate
	return am.authenticate(ctx)
}

// RefreshToken refreshes the access token
func (am *AuthManager) RefreshToken(ctx context.Context) error {
	am.authMutex.Lock()
	defer am.authMutex.Unlock()

	if am.auth == nil || am.auth.RefreshToken == "" {
		_, err := am.authenticateUnsafe(ctx)
		return err
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", am.auth.RefreshToken)
	data.Set("client_id", am.config.Credentials["client_id"])

	req := &common.HTTPRequest{
		Method: "POST",
		URL:    am.baseURL + "/authentification",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body:    []byte(data.Encode()),
		Timeout: am.config.Timeout,
	}

	resp, err := am.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh token request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token failed with status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(resp.Body, &authResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %w", err)
	}

	am.auth = &authResp
	am.logger.Debug("B-PAY token refreshed")

	return nil
}

// authenticate performs initial authentication
func (am *AuthManager) authenticate(ctx context.Context) (string, error) {
	am.authMutex.Lock()
	defer am.authMutex.Unlock()

	return am.authenticateUnsafe(ctx)
}

// authenticateUnsafe performs authentication without locking
func (am *AuthManager) authenticateUnsafe(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", am.config.Credentials["username"])
	data.Set("password", am.config.Credentials["password"])
	data.Set("client_id", am.config.Credentials["client_id"])

	req := &common.HTTPRequest{
		Method: "POST",
		URL:    am.baseURL + "/authentification",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body:    []byte(data.Encode()),
		Timeout: am.config.Timeout,
	}

	am.logger.Debug("Authenticating with B-PAY", "username", am.config.Credentials["username"])

	resp, err := am.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(resp.Body, &authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	am.auth = &authResp
	am.logger.Info("B-PAY authentication successful")

	return authResp.AccessToken, nil
}

// isTokenExpired checks if current token is expired
func (am *AuthManager) isTokenExpired() bool {
	if am.auth == nil {
		return true
	}

	// In a real implementation, you would track token issue time
	// and use expires_in to determine expiration
	// For simplicity, we assume token is valid for short time
	return false
}
