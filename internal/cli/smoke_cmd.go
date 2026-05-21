package cli

import (
	"time"

	"github.com/dannolan/ga4-cli/internal/ga4"
	"github.com/spf13/cobra"
)

func newSmokeCommand(ctx *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "smoke",
		Short: "Run a read-only authentication and report smoke test",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			end := time.Now().Format(time.DateOnly)
			start := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
			result, err := client.RunReport(cmd.Context(), ga4.ReportParams{
				PropertyID: property,
				Metrics:    []string{"sessions"},
				DateRanges: []ga4.DateRange{{StartDate: start, EndDate: end}},
				Limit:      1,
			})
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(map[string]any{
				"ok":          true,
				"property_id": property,
				"date_range":  map[string]string{"startDate": start, "endDate": end},
				"row_count":   result.RowCount,
			})
		},
	}
}
