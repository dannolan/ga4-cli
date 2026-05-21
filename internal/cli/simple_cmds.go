package cli

import (
	"github.com/spf13/cobra"
)

func newEventsCommand(ctx *appContext) *cobra.Command {
	var startDate, endDate, output string
	var limit int64
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Get top events report",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.TopEvents(cmd.Context(), property, defaultDateRange(startDate, endDate), limit)
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
	cmd.Flags().StringVarP(&startDate, "start-date", "s", "", "start date, YYYY-MM-DD")
	cmd.Flags().StringVarP(&endDate, "end-date", "e", "", "end date, YYYY-MM-DD")
	cmd.Flags().Int64VarP(&limit, "limit", "l", 20, "maximum rows")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	return cmd
}

func newPagesCommand(ctx *appContext) *cobra.Command {
	var startDate, endDate, output string
	var limit int64
	cmd := &cobra.Command{
		Use:   "pages",
		Short: "Get top pages report",
		RunE: func(cmd *cobra.Command, args []string) error {
			property, err := ctx.DefaultProperty()
			if err != nil {
				return err
			}
			client, err := ctx.Client(cmd.Context())
			if err != nil {
				return err
			}
			result, err := client.TopPages(cmd.Context(), property, defaultDateRange(startDate, endDate), limit)
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
	cmd.Flags().StringVarP(&startDate, "start-date", "s", "", "start date, YYYY-MM-DD")
	cmd.Flags().StringVarP(&endDate, "end-date", "e", "", "end date, YYYY-MM-DD")
	cmd.Flags().Int64VarP(&limit, "limit", "l", 20, "maximum rows")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	return cmd
}
