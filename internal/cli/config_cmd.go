package cli

import (
	"os"

	"github.com/dannolan/ga4-cli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect and write GA4 CLI config",
	}
	cmd.AddCommand(newConfigShowCommand(ctx), newConfigInitCommand(ctx))
	return cmd
}

func newConfigShowCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show non-secret config status",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ctx.Store()
			if err != nil {
				return err
			}
			cfg, err := ctx.Config()
			if err != nil {
				return err
			}
			legacyToken, _ := config.LegacyTokenPath()
			_, primaryErr := os.Stat(store.TokenPath())
			_, legacyErr := os.Stat(legacyToken)
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"config_dir":               store.Dir(),
				"client_id_configured":     cfg.ClientID != "",
				"client_secret_configured": cfg.ClientSecret != "",
				"property_id":              cfg.PropertyID,
				"token_configured":         primaryErr == nil || legacyErr == nil,
				"token_path":               store.TokenPath(),
				"legacy_token_path":        legacyToken,
			})
		},
	}
}

func newConfigInitCommand(ctx *appContext) *cobra.Command {
	var clientID, clientSecret, propertyID string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write config.json with OAuth client settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := ctx.Store()
			if err != nil {
				return err
			}
			return store.Save(config.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				PropertyID:   propertyID,
			})
		},
	}
	cmd.Flags().StringVar(&clientID, "client-id", "", "Google OAuth client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "Google OAuth client secret")
	cmd.Flags().StringVar(&propertyID, "property", "", "default GA4 property ID")
	_ = cmd.MarkFlagRequired("client-id")
	_ = cmd.MarkFlagRequired("client-secret")
	return cmd
}
