package cli

import (
	"github.com/dannolan/ga4-cli/internal/ga4"
	"github.com/spf13/cobra"
)

func newCompareCommand(ctx *appContext) *cobra.Command {
	var metrics, dimensions, currentStart, currentEnd, previousStart, previousEnd, output string
	var limit int64
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare GA4 data between two time periods",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			currentRange, previousRange := defaultComparisonRanges(currentStart, currentEnd, previousStart, previousEnd)
			base := ga4.ReportParams{
				PropertyID: property,
				Metrics:    splitCSV(metrics),
				Dimensions: splitCSV(dimensions),
				Limit:      limit,
			}
			currentParams := base
			currentParams.DateRanges = []ga4.DateRange{currentRange}
			previousParams := base
			previousParams.DateRanges = []ga4.DateRange{previousRange}
			current, err := client.RunReport(cmd.Context(), currentParams)
			if err != nil {
				return err
			}
			previous, err := client.RunReport(cmd.Context(), previousParams)
			if err != nil {
				return err
			}
			text, err := formatComparison(current, previous, ctx.Format)
			if err != nil {
				return err
			}
			return writeOutput(text, output)
		},
	}
	cmd.Flags().StringVarP(&metrics, "metrics", "m", "", "comma-separated metric names")
	cmd.Flags().StringVarP(&dimensions, "dimensions", "d", "", "comma-separated dimension names")
	cmd.Flags().StringVar(&currentStart, "current-start", "", "current period start date, YYYY-MM-DD")
	cmd.Flags().StringVar(&currentEnd, "current-end", "", "current period end date, YYYY-MM-DD")
	cmd.Flags().StringVar(&previousStart, "previous-start", "", "previous period start date, YYYY-MM-DD")
	cmd.Flags().StringVar(&previousEnd, "previous-end", "", "previous period end date, YYYY-MM-DD")
	cmd.Flags().Int64VarP(&limit, "limit", "l", 100, "maximum rows")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	_ = cmd.MarkFlagRequired("metrics")
	return cmd
}
