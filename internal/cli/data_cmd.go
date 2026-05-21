package cli

import (
	"github.com/spf13/cobra"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

func newDataCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Call read-only GA4 Data API methods",
	}
	cmd.AddCommand(
		newDataMetadataCommand(ctx),
		newDataRealtimeCommand(ctx),
		newDataCompatibilityCommand(ctx),
		newDataRawReportCommand(ctx),
		newDataBatchReportsCommand(ctx),
		newDataPivotCommand(ctx),
		newDataBatchPivotsCommand(ctx),
		newDataAudienceExportsCommand(ctx),
	)
	return cmd
}

func newDataMetadataCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata",
		Short: "Get GA4 dimensions and metrics metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.Metadata(cmd.Context(), property)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
}

func newDataRealtimeCommand(ctx *appContext) *cobra.Command {
	var metrics, dimensions string
	var limit int64
	cmd := &cobra.Command{
		Use:   "realtime",
		Short: "Run a realtime GA4 report",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			if metrics == "" {
				metrics = "activeUsers"
			}
			req := &analyticsdata.RunRealtimeReportRequest{
				Dimensions: dataDimensions(splitCSV(dimensions)),
				Metrics:    dataMetrics(splitCSV(metrics)),
				Limit:      limit,
			}
			result, err := client.RunRealtimeReport(cmd.Context(), property, req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVarP(&metrics, "metrics", "m", "activeUsers", "comma-separated realtime metric names")
	cmd.Flags().StringVarP(&dimensions, "dimensions", "d", "", "comma-separated realtime dimension names")
	cmd.Flags().Int64VarP(&limit, "limit", "l", 100, "maximum rows")
	return cmd
}

func newDataCompatibilityCommand(ctx *appContext) *cobra.Command {
	var metrics, dimensions, compatibility string
	cmd := &cobra.Command{
		Use:   "compatibility",
		Short: "Check dimension and metric compatibility",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			req := &analyticsdata.CheckCompatibilityRequest{
				Dimensions:          dataDimensions(splitCSV(dimensions)),
				Metrics:             dataMetrics(splitCSV(metrics)),
				CompatibilityFilter: compatibility,
			}
			result, err := client.CheckCompatibility(cmd.Context(), property, req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVarP(&metrics, "metrics", "m", "", "comma-separated metric names")
	cmd.Flags().StringVarP(&dimensions, "dimensions", "d", "", "comma-separated dimension names")
	cmd.Flags().StringVar(&compatibility, "compatibility", "", "compatibility filter, e.g. COMPATIBLE")
	return cmd
}

func newDataRawReportCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "run-report",
		Short: "Run properties.runReport with a JSON request body",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			var req analyticsdata.RunReportRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.RawRunReport(cmd.Context(), property, &req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "-", "JSON request body file, or - for stdin")
	return cmd
}

func newDataBatchReportsCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "batch-run-reports",
		Short: "Run properties.batchRunReports with a JSON request body",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			var req analyticsdata.BatchRunReportsRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.BatchRunReports(cmd.Context(), property, &req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "-", "JSON request body file, or - for stdin")
	return cmd
}

func newDataPivotCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "run-pivot-report",
		Short: "Run properties.runPivotReport with a JSON request body",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			var req analyticsdata.RunPivotReportRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.RunPivotReport(cmd.Context(), property, &req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "-", "JSON request body file, or - for stdin")
	return cmd
}

func newDataBatchPivotsCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "batch-run-pivot-reports",
		Short: "Run properties.batchRunPivotReports with a JSON request body",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			var req analyticsdata.BatchRunPivotReportsRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.BatchRunPivotReports(cmd.Context(), property, &req)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "-", "JSON request body file, or - for stdin")
	return cmd
}

func newDataAudienceExportsCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "audience-exports", Short: "Read GA4 audience exports"}
	var pageSize int64
	list := &cobra.Command{
		Use:   "list",
		Short: "List audience exports",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.AudienceExportsList(cmd.Context(), property, pageSize)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	list.Flags().Int64Var(&pageSize, "page-size", 50, "page size")
	get := &cobra.Command{
		Use:   "get <name>",
		Short: "Get an audience export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.AudienceExportGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	var limit, offset int64
	query := &cobra.Command{
		Use:   "query <name>",
		Short: "Query rows from an audience export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.AudienceExportQuery(cmd.Context(), args[0], &analyticsdata.QueryAudienceExportRequest{Limit: limit, Offset: offset})
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	query.Flags().Int64Var(&limit, "limit", 100, "maximum rows")
	query.Flags().Int64Var(&offset, "offset", 0, "row offset")
	cmd.AddCommand(list, get, query)
	return cmd
}

func dataDimensions(names []string) []*analyticsdata.Dimension {
	out := make([]*analyticsdata.Dimension, 0, len(names))
	for _, name := range names {
		out = append(out, &analyticsdata.Dimension{Name: name})
	}
	return out
}

func dataMetrics(names []string) []*analyticsdata.Metric {
	out := make([]*analyticsdata.Metric, 0, len(names))
	for _, name := range names {
		out = append(out, &analyticsdata.Metric{Name: name})
	}
	return out
}
