package cli

import (
	"github.com/spf13/cobra"
	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
)

func addAccountMutationCommands(ctx *appContext, parent *cobra.Command) {
	var patchOpts mutationOptions
	patch := &cobra.Command{
		Use:   "patch <account>",
		Short: "Patch an account from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaAccount
			if err := decodeJSONBody(patchOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, patchOpts.Apply, "accounts.patch", args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				call := service.Accounts.Patch(args[0], &req).Context(cmd.Context())
				if patchOpts.UpdateMask != "" {
					call.UpdateMask(patchOpts.UpdateMask)
				}
				return call.Do()
			})
		},
	}
	addPatchMutationFlags(patch, &patchOpts)

	var deleteApply bool
	deleteCmd := &cobra.Command{
		Use:   "delete <account>",
		Short: "Delete an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMutation(ctx, deleteApply, "accounts.delete", args[0], nil, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return service.Accounts.Delete(args[0]).Context(cmd.Context()).Do()
			})
		},
	}
	addApplyFlag(deleteCmd, &deleteApply)

	var provisionOpts mutationOptions
	provision := &cobra.Command{
		Use:   "provision-ticket",
		Short: "Provision an account ticket from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaProvisionAccountTicketRequest
			if err := decodeJSONBody(provisionOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, provisionOpts.Apply, "accounts.provisionAccountTicket", "", &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return service.Accounts.ProvisionAccountTicket(&req).Context(cmd.Context()).Do()
			})
		},
	}
	addMutationFlags(provision, &provisionOpts)
	parent.AddCommand(patch, deleteCmd, provision)
}

func addPropertyMutationCommands(ctx *appContext, parent *cobra.Command) {
	var createOpts mutationOptions
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a property from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaProperty
			if err := decodeJSONBody(createOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, createOpts.Apply, "properties.create", "", &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return service.Properties.Create(&req).Context(cmd.Context()).Do()
			})
		},
	}
	addMutationFlags(create, &createOpts)

	var patchOpts mutationOptions
	patch := &cobra.Command{
		Use:   "patch <property>",
		Short: "Patch a property from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaProperty
			if err := decodeJSONBody(patchOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, patchOpts.Apply, "properties.patch", args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				call := service.Properties.Patch(args[0], &req).Context(cmd.Context())
				if patchOpts.UpdateMask != "" {
					call.UpdateMask(patchOpts.UpdateMask)
				}
				return call.Do()
			})
		},
	}
	addPatchMutationFlags(patch, &patchOpts)

	var deleteApply bool
	deleteCmd := &cobra.Command{
		Use:   "delete <property>",
		Short: "Delete a property",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMutation(ctx, deleteApply, "properties.delete", args[0], nil, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return service.Properties.Delete(args[0]).Context(cmd.Context()).Do()
			})
		},
	}
	addApplyFlag(deleteCmd, &deleteApply)

	var acknowledgeOpts mutationOptions
	acknowledge := &cobra.Command{
		Use:   "acknowledge-user-data-collection <property>",
		Short: "Acknowledge user data collection from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaAcknowledgeUserDataCollectionRequest
			if err := decodeJSONBody(acknowledgeOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, acknowledgeOpts.Apply, "properties.acknowledgeUserDataCollection", args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return service.Properties.AcknowledgeUserDataCollection(args[0], &req).Context(cmd.Context()).Do()
			})
		},
	}
	addMutationFlags(acknowledge, &acknowledgeOpts)

	var retentionOpts mutationOptions
	retention := &cobra.Command{
		Use:   "update-data-retention-settings <data-retention-settings>",
		Short: "Update data retention settings from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaDataRetentionSettings
			if err := decodeJSONBody(retentionOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, retentionOpts.Apply, "properties.updateDataRetentionSettings", args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				call := service.Properties.UpdateDataRetentionSettings(args[0], &req).Context(cmd.Context())
				if retentionOpts.UpdateMask != "" {
					call.UpdateMask(retentionOpts.UpdateMask)
				}
				return call.Do()
			})
		},
	}
	addPatchMutationFlags(retention, &retentionOpts)
	parent.AddCommand(create, patch, deleteCmd, acknowledge, retention)
}
