package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func NewRootCommand() *cobra.Command {
	ctx := &appContext{}
	root := &cobra.Command{
		Use:           "ga4",
		Short:         "Agent-first CLI for Google Analytics 4",
		Long:          "ga4 is a compiled Go CLI for Google Analytics 4 reports. It reuses local OAuth credentials and emits stable table, JSON, markdown, or CSV output.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&ctx.ConfigDir, "config-dir", "", "config directory; defaults to ~/.config/ga4-cli")
	root.PersistentFlags().StringVarP(&ctx.Property, "property", "p", "", "GA4 property ID; defaults to GA4_PROPERTY_ID")
	root.PersistentFlags().StringVarP(&ctx.Format, "format", "f", "table", "output format: table, json, markdown, csv")
	root.PersistentFlags().BoolVar(&ctx.JSON, "json", false, "emit agent-readable JSON for metadata commands")
	root.PersistentFlags().BoolVar(&ctx.Pretty, "pretty", true, "pretty-print JSON")

	root.AddCommand(
		newAuthCommand(ctx),
		newDataCommand(ctx),
		newAdminCommand(ctx),
		newReportCommand(ctx),
		newCompareCommand(ctx),
		newEventsCommand(ctx),
		newPagesCommand(ctx),
		newConfigCommand(ctx),
		newSmokeCommand(ctx),
		newDoctorCommand(ctx),
		newVersionCommand(ctx),
		newManifestCommand(ctx, root),
	)
	return root
}

func newVersionCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"version": Version,
				"commit":  Commit,
				"date":    Date,
			})
		},
	}
}

func newManifestCommand(ctx *appContext, root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "manifest",
		Short: "Print an agent-readable command manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			commands := []map[string]any{}
			var walk func(*cobra.Command)
			walk = func(parent *cobra.Command) {
				for _, child := range parent.Commands() {
					if child.Hidden {
						continue
					}
					flags := []string{}
					child.Flags().VisitAll(func(flag *pflag.Flag) {
						flags = append(flags, "--"+flag.Name)
					})
					commands = append(commands, map[string]any{
						"path":     child.CommandPath(),
						"use":      child.UseLine(),
						"short":    child.Short,
						"mutation": false,
						"flags":    flags,
					})
					walk(child)
				}
			}
			walk(root)
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"name":     "ga4-cli",
				"binary":   "ga4",
				"version":  Version,
				"commands": commands,
				"policy": map[string]any{
					"json":        "Use --format json for report data or --json for metadata commands.",
					"credentials": "Reads CLIENT_ID, CLIENT_SECRET, and GA4_PROPERTY_ID from env or ~/.config/ga4-cli/env.",
					"token":       "Reuses OAuth token from ~/.config/ga4-cli/token.json or legacy ~/.ga4-cli/token.json.",
				},
			})
		},
	}
}
