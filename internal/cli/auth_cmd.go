package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dannolan/ga4-cli/internal/config"
	"github.com/dannolan/ga4-cli/internal/ga4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func newAuthCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Google Analytics OAuth tokens",
	}
	cmd.AddCommand(newAuthStatusCommand(ctx), newAuthLoginCommand(ctx), newAuthLogoutCommand(ctx))
	return cmd
}

func newAuthStatusCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show OAuth token status",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ctx.Store()
			if err != nil {
				return err
			}
			legacyPath, _ := config.LegacyTokenPath()
			token, err := ga4.LoadToken(ga4.TokenStore{PrimaryPath: store.TokenPath(), LegacyPath: legacyPath})
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"token_configured":  token != nil,
				"token_valid":       token != nil && token.Valid(),
				"expiry":            expiryString(token),
				"token_path":        store.TokenPath(),
				"legacy_token_path": legacyPath,
			})
		},
	}
}

func newAuthLoginCommand(ctx *appContext) *cobra.Command {
	var noBrowser bool
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authorize this CLI with Google Analytics",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := ctx.Config()
			if err != nil {
				return err
			}
			store, err := ctx.Store()
			if err != nil {
				return err
			}
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return err
			}
			defer listener.Close()

			redirectURL := "http://" + listener.Addr().String() + "/oauth2callback"
			state, err := ga4.RandomState()
			if err != nil {
				return err
			}
			authURL := ga4.AuthCodeURL(cfg, redirectURL, state)
			fmt.Fprintln(os.Stderr, "Open this URL to authorize GA4 CLI:")
			fmt.Fprintln(os.Stderr, authURL)
			if !noBrowser {
				_ = openBrowser(authURL)
			}

			codeCh := make(chan string, 1)
			errCh := make(chan error, 1)
			server := &http.Server{ReadHeaderTimeout: 5 * time.Second}
			server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/oauth2callback" {
					http.NotFound(w, r)
					return
				}
				if r.URL.Query().Get("state") != state {
					http.Error(w, "state mismatch", http.StatusBadRequest)
					errCh <- fmt.Errorf("OAuth state mismatch")
					return
				}
				if errText := r.URL.Query().Get("error"); errText != "" {
					http.Error(w, errText, http.StatusBadRequest)
					errCh <- fmt.Errorf("OAuth error: %s", errText)
					return
				}
				code := r.URL.Query().Get("code")
				if code == "" {
					http.Error(w, "missing code", http.StatusBadRequest)
					errCh <- fmt.Errorf("OAuth callback did not include a code")
					return
				}
				fmt.Fprintln(w, "GA4 CLI authorization complete. You can close this tab.")
				codeCh <- code
			})
			go func() {
				if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
					errCh <- err
				}
			}()
			defer server.Shutdown(context.Background())

			var code string
			select {
			case code = <-codeCh:
			case err := <-errCh:
				return err
			case <-time.After(2 * time.Minute):
				return fmt.Errorf("timed out waiting for OAuth callback")
			}
			token, err := ga4.ExchangeCode(cmd.Context(), cfg, redirectURL, code)
			if err != nil {
				return err
			}
			if err := ga4.SaveToken(store.TokenPath(), token); err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"ok":         true,
				"token_path": store.TokenPath(),
				"expiry":     expiryString(token),
			})
		},
	}
	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "print auth URL without opening a browser")
	return cmd
}

func newAuthLogoutCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the local OAuth token",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ctx.Store()
			if err != nil {
				return err
			}
			err = os.Remove(store.TokenPath())
			if os.IsNotExist(err) {
				err = nil
			}
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(map[string]any{"ok": true, "removed": store.TokenPath()})
		},
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func expiryString(token *oauth2.Token) string {
	if token == nil || token.Expiry.IsZero() {
		return ""
	}
	return token.Expiry.Format(time.RFC3339)
}
