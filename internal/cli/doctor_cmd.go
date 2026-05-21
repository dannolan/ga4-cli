package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dannolan/ga4-cli/internal/config"
	"github.com/dannolan/ga4-cli/internal/ga4"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/googleapi"

	"github.com/spf13/cobra"
)

func newDoctorCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check local GA4 CLI readiness",
		RunE: func(cmd *cobra.Command, args []string) error {
			report := map[string]any{
				"ok":      true,
				"version": Version,
				"commit":  Commit,
				"checks":  map[string]any{},
			}
			checks := report["checks"].(map[string]any)

			store, err := ctx.Store()
			if err != nil {
				addDoctorCheck(checks, "config_store", false, err, nil)
				report["ok"] = false
				ctx.JSON = true
				return ctx.Print(report)
			}
			addDoctorCheck(checks, "config_store", true, nil, map[string]any{"config_dir": store.Dir()})

			cfg, err := ctx.Config()
			if err != nil {
				addDoctorCheck(checks, "credentials", false, err, nil)
				report["ok"] = false
			} else {
				addDoctorCheck(checks, "credentials", true, nil, map[string]any{
					"client_id_configured":     cfg.ClientID != "",
					"client_secret_configured": cfg.ClientSecret != "",
					"property_id":              cfg.PropertyID,
				})
			}

			legacyPath, _ := config.LegacyTokenPath()
			token, tokenErr := ga4.LoadToken(ga4.TokenStore{PrimaryPath: store.TokenPath(), LegacyPath: legacyPath})
			if tokenErr != nil {
				addDoctorCheck(checks, "oauth_token", false, tokenErr, nil)
				report["ok"] = false
			} else {
				addDoctorCheck(checks, "oauth_token", token != nil, nil, map[string]any{
					"token_configured":  token != nil,
					"token_valid":       token != nil && token.Valid(),
					"expiry":            expiryString(token),
					"token_path":        store.TokenPath(),
					"legacy_token_path": legacyPath,
				})
				if token == nil {
					report["ok"] = false
				}
			}

			property := ""
			if cfg.PropertyID != "" {
				property = cfg.PropertyID
			}
			if property == "" {
				addDoctorCheck(checks, "property", false, errors.New("property ID is required; pass --property or set GA4_PROPERTY_ID"), nil)
				report["ok"] = false
			} else {
				addDoctorCheck(checks, "property", true, nil, map[string]any{"property_id": property})
			}

			if err == nil && token != nil && property != "" {
				client, clientErr := ctx.Client(cmd.Context())
				if clientErr != nil {
					addDoctorCheck(checks, "data_api", false, clientErr, nil)
					report["ok"] = false
				} else {
					end := time.Now().Format(time.DateOnly)
					start := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
					result, dataErr := client.RawRunReport(cmd.Context(), property, &analyticsdata.RunReportRequest{
						DateRanges: []*analyticsdata.DateRange{{StartDate: start, EndDate: end}},
						Metrics:    []*analyticsdata.Metric{{Name: "sessions"}},
						Limit:      1,
					})
					addDoctorCheck(checks, "data_api", dataErr == nil, dataErr, map[string]any{
						"row_count":  valueOrZero(result),
						"date_range": map[string]string{"startDate": start, "endDate": end},
					})
					if dataErr != nil {
						report["ok"] = false
					}
				}

				admin, adminErr := ctx.AdminService(cmd.Context())
				if adminErr != nil {
					addDoctorCheck(checks, "admin_api", false, adminErr, nil)
					report["ok"] = false
				} else {
					_, adminErr = admin.AccountSummaries.List().PageSize(1).Context(cmd.Context()).Do()
					detail := map[string]any{}
					if serviceDisabled(adminErr) {
						detail["enable_url"] = "https://console.developers.google.com/apis/api/analyticsadmin.googleapis.com/overview?project=1013065713343"
						detail["blocked_by"] = "analyticsadmin.googleapis.com is disabled for the OAuth project"
					}
					addDoctorCheck(checks, "admin_api", adminErr == nil, adminErr, detail)
					if adminErr != nil {
						report["ok"] = false
					}
				}
			}

			ctx.JSON = true
			return ctx.Print(report)
		},
	}
}

func addDoctorCheck(checks map[string]any, name string, ok bool, err error, details map[string]any) {
	check := map[string]any{"ok": ok}
	if err != nil {
		check["error"] = err.Error()
	}
	for key, value := range details {
		check[key] = value
	}
	checks[name] = check
}

func valueOrZero(result *analyticsdata.RunReportResponse) int64 {
	if result == nil {
		return 0
	}
	return result.RowCount
}

func serviceDisabled(err error) bool {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		for _, detail := range apiErr.Details {
			text := strings.ToLower(fmt.Sprint(detail))
			if strings.Contains(text, "service_disabled") || strings.Contains(text, "analyticsadmin.googleapis.com") {
				return true
			}
		}
	}
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "analyticsadmin.googleapis.com")
}
