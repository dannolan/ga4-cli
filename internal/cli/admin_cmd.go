package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
)

func newAdminCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Call read-only Google Analytics Admin API methods",
	}
	cmd.AddCommand(
		newAdminAccountSummariesCommand(ctx),
		newAdminAccountsCommand(ctx),
		newAdminPropertiesCommand(ctx),
		newAdminPropertyResourcesCommand(ctx),
	)
	return cmd
}

func newAdminAccountSummariesCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "account-summaries", Short: "Read account summaries"}
	var pageSize int64
	list := &cobra.Command{
		Use:   "list",
		Short: "List account summaries",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			call := service.AccountSummaries.List().Context(cmd.Context())
			if pageSize > 0 {
				call.PageSize(pageSize)
			}
			result, err := call.Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	list.Flags().Int64Var(&pageSize, "page-size", 50, "page size")
	cmd.AddCommand(list)
	return cmd
}

func newAdminAccountsCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "accounts", Short: "Read Analytics accounts"}
	var pageSize int64
	list := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			call := service.Accounts.List().Context(cmd.Context())
			if pageSize > 0 {
				call.PageSize(pageSize)
			}
			result, err := call.Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	list.Flags().Int64Var(&pageSize, "page-size", 50, "page size")
	get := &cobra.Command{
		Use:   "get <account>",
		Short: "Get an account, e.g. accounts/123",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Accounts.Get(args[0]).Context(cmd.Context()).Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	sharing := &cobra.Command{
		Use:   "data-sharing-settings <account>",
		Short: "Get account data sharing settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Accounts.GetDataSharingSettings(args[0] + "/dataSharingSettings").Context(cmd.Context()).Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.AddCommand(list, get, sharing, newAdminAccountAccessReportCommand(ctx), newAdminChangeHistoryCommand(ctx))
	addAccountMutationCommands(ctx, cmd)
	return cmd
}

func newAdminPropertiesCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "properties", Short: "Read Analytics properties"}
	var pageSize int64
	var filter string
	list := &cobra.Command{
		Use:   "list",
		Short: "List properties",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			call := service.Properties.List().Context(cmd.Context())
			if pageSize > 0 {
				call.PageSize(pageSize)
			}
			if filter != "" {
				call.Filter(filter)
			}
			result, err := call.Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	list.Flags().Int64Var(&pageSize, "page-size", 50, "page size")
	list.Flags().StringVar(&filter, "filter", "", `property filter, e.g. parent:accounts/123`)
	get := &cobra.Command{
		Use:   "get <property>",
		Short: "Get a property, e.g. properties/123",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Properties.Get(args[0]).Context(cmd.Context()).Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	retention := &cobra.Command{
		Use:   "data-retention-settings <property>",
		Short: "Get property data retention settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Properties.GetDataRetentionSettings(args[0] + "/dataRetentionSettings").Context(cmd.Context()).Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.AddCommand(list, get, retention, newAdminPropertyAccessReportCommand(ctx))
	addPropertyMutationCommands(ctx, cmd)
	return cmd
}

func newAdminAccountAccessReportCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "access-report <account>",
		Short: "Run an account access report with a JSON request body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaRunAccessReportRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Accounts.RunAccessReport(args[0], &req).Context(cmd.Context()).Do()
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

func newAdminPropertyAccessReportCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "access-report <property>",
		Short: "Run a property access report with a JSON request body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaRunAccessReportRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Properties.RunAccessReport(args[0], &req).Context(cmd.Context()).Do()
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

func newAdminChangeHistoryCommand(ctx *appContext) *cobra.Command {
	var body string
	cmd := &cobra.Command{
		Use:   "change-history <account>",
		Short: "Search account change history with a JSON request body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req analyticsadmin.GoogleAnalyticsAdminV1betaSearchChangeHistoryEventsRequest
			if err := decodeJSONBody(body, &req); err != nil {
				return err
			}
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Accounts.SearchChangeHistoryEvents(args[0], &req).Context(cmd.Context()).Do()
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

func newAdminPropertyResourcesCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "property-resources", Short: "Manage property child resources"}
	conversionEvents := propertyListGetCommand(ctx, "conversion-events", "conversion event", listConversionEvents, getConversionEvent)
	customDimensions := propertyListGetCommand(ctx, "custom-dimensions", "custom dimension", listCustomDimensions, getCustomDimension)
	customMetrics := propertyListGetCommand(ctx, "custom-metrics", "custom metric", listCustomMetrics, getCustomMetric)
	dataStreams := propertyListGetCommand(ctx, "data-streams", "data stream", listDataStreams, getDataStream)
	firebaseLinks := propertyListOnlyCommand(ctx, "firebase-links", "Firebase link", listFirebaseLinks)
	googleAdsLinks := propertyListOnlyCommand(ctx, "google-ads-links", "Google Ads link", listGoogleAdsLinks)
	keyEvents := propertyListGetCommand(ctx, "key-events", "key event", listKeyEvents, getKeyEvent)
	measurementSecrets := measurementSecretsCommand(ctx)
	addConversionEventMutationCommands(ctx, conversionEvents)
	addCustomDimensionMutationCommands(ctx, customDimensions)
	addCustomMetricMutationCommands(ctx, customMetrics)
	addDataStreamMutationCommands(ctx, dataStreams)
	addFirebaseLinkMutationCommands(ctx, firebaseLinks)
	addGoogleAdsLinkMutationCommands(ctx, googleAdsLinks)
	addKeyEventMutationCommands(ctx, keyEvents)
	addMeasurementSecretMutationCommands(ctx, measurementSecrets)
	cmd.AddCommand(conversionEvents, customDimensions, customMetrics, dataStreams, firebaseLinks, googleAdsLinks, keyEvents, measurementSecrets)
	return cmd
}

type listFunc func(*analyticsadmin.Service, string, int64) (any, error)
type getFunc func(*analyticsadmin.Service, string) (any, error)

func propertyListGetCommand(ctx *appContext, use, label string, list listFunc, get getFunc) *cobra.Command {
	cmd := propertyListOnlyCommand(ctx, use, label, list)
	cmd.AddCommand(&cobra.Command{
		Use:   "get <name>",
		Short: fmt.Sprintf("Get a %s", label),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := get(service, args[0])
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	})
	return cmd
}

func propertyListOnlyCommand(ctx *appContext, use, label string, list listFunc) *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: "Read " + label + " resources"}
	var pageSize int64
	listCmd := &cobra.Command{
		Use:   "list <property>",
		Short: fmt.Sprintf("List %s resources", label),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := list(service, args[0], pageSize)
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	listCmd.Flags().Int64Var(&pageSize, "page-size", 50, "page size")
	cmd.AddCommand(listCmd)
	return cmd
}

func measurementSecretsCommand(ctx *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "measurement-protocol-secrets", Short: "Read measurement protocol secrets"}
	var pageSize int64
	list := &cobra.Command{
		Use:   "list <data-stream>",
		Short: "List measurement protocol secrets for a data stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			call := service.Properties.DataStreams.MeasurementProtocolSecrets.List(args[0]).Context(cmd.Context())
			if pageSize > 0 {
				call.PageSize(pageSize)
			}
			result, err := call.Do()
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
		Short: "Get a measurement protocol secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := ctx.AdminService(cmd.Context())
			if err != nil {
				return err
			}
			result, err := service.Properties.DataStreams.MeasurementProtocolSecrets.Get(args[0]).Context(cmd.Context()).Do()
			if err != nil {
				return err
			}
			ctx.JSON = true
			return ctx.Print(result)
		},
	}
	cmd.AddCommand(list, get)
	return cmd
}

func listConversionEvents(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	call := s.Properties.ConversionEvents.List(parent)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}
func getConversionEvent(s *analyticsadmin.Service, name string) (any, error) {
	return s.Properties.ConversionEvents.Get(name).Do()
}
func listCustomDimensions(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	call := s.Properties.CustomDimensions.List(parent)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}
func getCustomDimension(s *analyticsadmin.Service, name string) (any, error) {
	return s.Properties.CustomDimensions.Get(name).Do()
}
func listCustomMetrics(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	call := s.Properties.CustomMetrics.List(parent)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}
func getCustomMetric(s *analyticsadmin.Service, name string) (any, error) {
	return s.Properties.CustomMetrics.Get(name).Do()
}
func listDataStreams(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	call := s.Properties.DataStreams.List(parent)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}
func getDataStream(s *analyticsadmin.Service, name string) (any, error) {
	return s.Properties.DataStreams.Get(name).Do()
}
func listFirebaseLinks(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	return s.Properties.FirebaseLinks.List(parent).Do()
}
func listGoogleAdsLinks(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	return s.Properties.GoogleAdsLinks.List(parent).Do()
}
func listKeyEvents(s *analyticsadmin.Service, parent string, pageSize int64) (any, error) {
	call := s.Properties.KeyEvents.List(parent)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}
func getKeyEvent(s *analyticsadmin.Service, name string) (any, error) {
	return s.Properties.KeyEvents.Get(name).Do()
}
