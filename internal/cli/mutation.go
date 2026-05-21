package cli

import "github.com/spf13/cobra"

type mutationOptions struct {
	Body       string
	UpdateMask string
	Apply      bool
}

func addMutationFlags(cmd *cobra.Command, opts *mutationOptions) {
	cmd.Flags().StringVar(&opts.Body, "body", "-", "JSON request body file, or - for stdin")
	cmd.Flags().BoolVar(&opts.Apply, "apply", false, "apply the mutation; without this flag the command only prints a dry run")
}

func addPatchMutationFlags(cmd *cobra.Command, opts *mutationOptions) {
	addMutationFlags(cmd, opts)
	cmd.Flags().StringVar(&opts.UpdateMask, "update-mask", "", "comma-separated field mask for patch/update calls")
}

func addApplyFlag(cmd *cobra.Command, apply *bool) {
	cmd.Flags().BoolVar(apply, "apply", false, "apply the mutation; without this flag the command only prints a dry run")
}

func printDryRun(ctx *appContext, operation string, target string, request any) error {
	ctx.JSON = true
	return ctx.Print(map[string]any{
		"dry_run":    true,
		"operation":  operation,
		"target":     target,
		"request":    request,
		"apply_hint": "pass --apply to call the Google API",
	})
}

func runMutation(ctx *appContext, apply bool, operation string, target string, request any, run func() (any, error)) error {
	if !apply {
		return printDryRun(ctx, operation, target, request)
	}
	result, err := run()
	if err != nil {
		return err
	}
	ctx.JSON = true
	return ctx.Print(map[string]any{
		"applied":   true,
		"operation": operation,
		"target":    target,
		"result":    result,
	})
}
