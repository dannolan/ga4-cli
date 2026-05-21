package ga4

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dannolan/ga4-cli/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	analyticsReadonlyScope = "https://www.googleapis.com/auth/analytics.readonly"
	analyticsEditScope     = "https://www.googleapis.com/auth/analytics.edit"
)

type TokenStore struct {
	PrimaryPath string
	LegacyPath  string
}

func OAuthConfig(cfg config.Config, redirectURL string) *oauth2.Config {
	if redirectURL == "" {
		redirectURL = "urn:ietf:wg:oauth:2.0:oob"
	}
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{analyticsReadonlyScope, analyticsEditScope},
		Endpoint:     google.Endpoint,
	}
}

func AuthenticatedHTTPClient(ctx context.Context, cfg config.Config, store TokenStore) (*http.Client, error) {
	token, err := LoadToken(store)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("no OAuth token found; run the existing Node auth once or place token.json at %s", store.PrimaryPath)
	}
	source := OAuthConfig(cfg, "").TokenSource(ctx, token)
	refreshed, err := source.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh access token: %w", err)
	}
	if refreshed.AccessToken != token.AccessToken {
		_ = SaveToken(store.PrimaryPath, refreshed)
	}
	return oauth2.NewClient(ctx, source), nil
}

func AuthCodeURL(cfg config.Config, redirectURL string, state string) string {
	return OAuthConfig(cfg, redirectURL).AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func ExchangeCode(ctx context.Context, cfg config.Config, redirectURL string, code string) (*oauth2.Token, error) {
	return OAuthConfig(cfg, redirectURL).Exchange(ctx, code)
}

func RandomState() (string, error) {
	var raw [24]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func LoadToken(store TokenStore) (*oauth2.Token, error) {
	for _, path := range []string{store.PrimaryPath, store.LegacyPath} {
		if path == "" {
			continue
		}
		data, err := os.ReadFile(path)
		if err == nil {
			var token oauth2.Token
			if err := json.Unmarshal(data, &token); err != nil {
				return nil, fmt.Errorf("read token %s: %w", path, err)
			}
			return &token, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, nil
}

func SaveToken(path string, token *oauth2.Token) error {
	if path == "" || token == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}
