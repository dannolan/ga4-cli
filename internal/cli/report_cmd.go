package cli

import (
	"strings"

	"github.com/dannolan/ga4-cli/internal/ga4"
	"github.com/spf13/cobra"
)

func newReportCommand(ctx *appContext) *cobra.Command {
	var metrics, dimensions, startDate, endDate, output, orderBy string
	var limit, offset int64
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a custom GA4 report",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.RunReport(cmd.Context(), ga4.ReportParams{
				PropertyID: property,
				Metrics:    splitCSV(metrics),
				Dimensions: splitCSV(dimensions),
				DateRanges: []ga4.DateRange{defaultDateRange(startDate, endDate)},
				Limit:      limit,
				Offset:     offset,
				OrderBys:   parseOrderBys(orderBy),
			})
			if err != nil {
				return err
			}
			text, err := formatReport(result, ctx.Format)
			if err != nil {
				return err
			}
			return writeOutput(text, output)
		},
	}
	cmd.Flags().StringVarP(&metrics, "metrics", "m", "", "comma-separated metric names")
	cmd.Flags().StringVarP(&dimensions, "dimensions", "d", "", "comma-separated dimension names")
	cmd.Flags().StringVarP(&startDate, "start-date", "s", "", "start date, YYYY-MM-DD")
	cmd.Flags().StringVarP(&endDate, "end-date", "e", "", "end date, YYYY-MM-DD")
	cmd.Flags().Int64VarP(&limit, "limit", "l", 100, "maximum rows")
	cmd.Flags().Int64Var(&offset, "offset", 0, "row offset")
	cmd.Flags().StringVar(&orderBy, "order-by", "", "comma-separated metric names to order by; prefix with - for descending")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	_ = cmd.MarkFlagRequired("metrics")
	return cmd
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseOrderBys(value string) []ga4.OrderBy {
	parts := splitCSV(value)
	out := make([]ga4.OrderBy, 0, len(parts))
	for _, part := range parts {
		desc := strings.HasPrefix(part, "-")
		out = append(out, ga4.OrderBy{Field: strings.TrimPrefix(part, "-"), Desc: desc})
	}
	return out
}
