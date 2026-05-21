package cli

import (
	"github.com/spf13/cobra"
	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
)

func addConversionEventMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaConversionEvent]{
		CreateUse: "create <property>", CreateOperation: "properties.conversionEvents.create",
		PatchUse: "patch <conversion-event>", PatchOperation: "properties.conversionEvents.patch",
		DeleteUse: "delete <conversion-event>", DeleteOperation: "properties.conversionEvents.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaConversionEvent) (any, error) {
			return s.Properties.ConversionEvents.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaConversionEvent, mask string) (any, error) {
			call := s.Properties.ConversionEvents.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.ConversionEvents.Delete(name).Do()
		},
	})
}

func addCustomDimensionMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchArchive(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaCustomDimension]{
		CreateUse: "create <property>", CreateOperation: "properties.customDimensions.create",
		PatchUse: "patch <custom-dimension>", PatchOperation: "properties.customDimensions.patch",
		DeleteUse: "archive <custom-dimension>", DeleteOperation: "properties.customDimensions.archive",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaCustomDimension) (any, error) {
			return s.Properties.CustomDimensions.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaCustomDimension, mask string) (any, error) {
			call := s.Properties.CustomDimensions.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		ArchiveDimension: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.CustomDimensions.Archive(name, &analyticsadmin.GoogleAnalyticsAdminV1betaArchiveCustomDimensionRequest{}).Do()
		},
	})
}

func addCustomMetricMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchArchive(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaCustomMetric]{
		CreateUse: "create <property>", CreateOperation: "properties.customMetrics.create",
		PatchUse: "patch <custom-metric>", PatchOperation: "properties.customMetrics.patch",
		DeleteUse: "archive <custom-metric>", DeleteOperation: "properties.customMetrics.archive",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaCustomMetric) (any, error) {
			return s.Properties.CustomMetrics.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaCustomMetric, mask string) (any, error) {
			call := s.Properties.CustomMetrics.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		ArchiveMetric: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.CustomMetrics.Archive(name, &analyticsadmin.GoogleAnalyticsAdminV1betaArchiveCustomMetricRequest{}).Do()
		},
	})
}

func addDataStreamMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaDataStream]{
		CreateUse: "create <property>", CreateOperation: "properties.dataStreams.create",
		PatchUse: "patch <data-stream>", PatchOperation: "properties.dataStreams.patch",
		DeleteUse: "delete <data-stream>", DeleteOperation: "properties.dataStreams.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaDataStream) (any, error) {
			return s.Properties.DataStreams.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaDataStream, mask string) (any, error) {
			call := s.Properties.DataStreams.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.DataStreams.Delete(name).Do()
		},
	})
}

func addMeasurementSecretMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaMeasurementProtocolSecret]{
		CreateUse: "create <data-stream>", CreateOperation: "properties.dataStreams.measurementProtocolSecrets.create",
		PatchUse: "patch <measurement-protocol-secret>", PatchOperation: "properties.dataStreams.measurementProtocolSecrets.patch",
		DeleteUse: "delete <measurement-protocol-secret>", DeleteOperation: "properties.dataStreams.measurementProtocolSecrets.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaMeasurementProtocolSecret) (any, error) {
			return s.Properties.DataStreams.MeasurementProtocolSecrets.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaMeasurementProtocolSecret, mask string) (any, error) {
			call := s.Properties.DataStreams.MeasurementProtocolSecrets.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.DataStreams.MeasurementProtocolSecrets.Delete(name).Do()
		},
	})
}

func addFirebaseLinkMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreateDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaFirebaseLink]{
		CreateUse: "create <property>", CreateOperation: "properties.firebaseLinks.create",
		DeleteUse: "delete <firebase-link>", DeleteOperation: "properties.firebaseLinks.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaFirebaseLink) (any, error) {
			return s.Properties.FirebaseLinks.Create(parent, req).Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.FirebaseLinks.Delete(name).Do()
		},
	})
}

func addGoogleAdsLinkMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaGoogleAdsLink]{
		CreateUse: "create <property>", CreateOperation: "properties.googleAdsLinks.create",
		PatchUse: "patch <google-ads-link>", PatchOperation: "properties.googleAdsLinks.patch",
		DeleteUse: "delete <google-ads-link>", DeleteOperation: "properties.googleAdsLinks.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaGoogleAdsLink) (any, error) {
			return s.Properties.GoogleAdsLinks.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaGoogleAdsLink, mask string) (any, error) {
			call := s.Properties.GoogleAdsLinks.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.GoogleAdsLinks.Delete(name).Do()
		},
	})
}

func addKeyEventMutationCommands(ctx *appContext, parent *cobra.Command) {
	addCreatePatchDelete(ctx, parent, resourceMutationSpec[analyticsadmin.GoogleAnalyticsAdminV1betaKeyEvent]{
		CreateUse: "create <property>", CreateOperation: "properties.keyEvents.create",
		PatchUse: "patch <key-event>", PatchOperation: "properties.keyEvents.patch",
		DeleteUse: "delete <key-event>", DeleteOperation: "properties.keyEvents.delete",
		Create: func(s *analyticsadmin.Service, parent string, req *analyticsadmin.GoogleAnalyticsAdminV1betaKeyEvent) (any, error) {
			return s.Properties.KeyEvents.Create(parent, req).Do()
		},
		Patch: func(s *analyticsadmin.Service, name string, req *analyticsadmin.GoogleAnalyticsAdminV1betaKeyEvent, mask string) (any, error) {
			call := s.Properties.KeyEvents.Patch(name, req)
			if mask != "" {
				call.UpdateMask(mask)
			}
			return call.Do()
		},
		Delete: func(s *analyticsadmin.Service, name string) (any, error) {
			return s.Properties.KeyEvents.Delete(name).Do()
		},
	})
}

type resourceMutationSpec[T any] struct {
	CreateUse, CreateOperation string
	PatchUse, PatchOperation   string
	DeleteUse, DeleteOperation string
	Create                     func(*analyticsadmin.Service, string, *T) (any, error)
	Patch                      func(*analyticsadmin.Service, string, *T, string) (any, error)
	Delete                     func(*analyticsadmin.Service, string) (any, error)
	ArchiveDimension           func(*analyticsadmin.Service, string) (any, error)
	ArchiveMetric              func(*analyticsadmin.Service, string) (any, error)
}

func addCreatePatchDelete[T any](ctx *appContext, parent *cobra.Command, spec resourceMutationSpec[T]) {
	addCreateDelete(ctx, parent, spec)
	var opts mutationOptions
	cmd := &cobra.Command{
		Use:   spec.PatchUse,
		Short: "Patch from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req T
			if err := decodeJSONBody(opts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, opts.Apply, spec.PatchOperation, args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return spec.Patch(service, args[0], &req, opts.UpdateMask)
			})
		},
	}
	addPatchMutationFlags(cmd, &opts)
	parent.AddCommand(cmd)
}

func addCreateDelete[T any](ctx *appContext, parent *cobra.Command, spec resourceMutationSpec[T]) {
	var createOpts mutationOptions
	create := &cobra.Command{
		Use:   spec.CreateUse,
		Short: "Create from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req T
			if err := decodeJSONBody(createOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, createOpts.Apply, spec.CreateOperation, args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return spec.Create(service, args[0], &req)
			})
		},
	}
	addMutationFlags(create, &createOpts)
	parent.AddCommand(create)

	var deleteApply bool
	deleteCmd := &cobra.Command{
		Use:   spec.DeleteUse,
		Short: "Delete the resource",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMutation(ctx, deleteApply, spec.DeleteOperation, args[0], nil, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return spec.Delete(service, args[0])
			})
		},
	}
	addApplyFlag(deleteCmd, &deleteApply)
	parent.AddCommand(deleteCmd)
}

func addCreatePatchArchive[T any](ctx *appContext, parent *cobra.Command, spec resourceMutationSpec[T]) {
	var createOpts mutationOptions
	create := &cobra.Command{
		Use:   spec.CreateUse,
		Short: "Create from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req T
			if err := decodeJSONBody(createOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, createOpts.Apply, spec.CreateOperation, args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return spec.Create(service, args[0], &req)
			})
		},
	}
	addMutationFlags(create, &createOpts)

	var patchOpts mutationOptions
	patch := &cobra.Command{
		Use:   spec.PatchUse,
		Short: "Patch from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var req T
			if err := decodeJSONBody(patchOpts.Body, &req); err != nil {
				return err
			}
			return runMutation(ctx, patchOpts.Apply, spec.PatchOperation, args[0], &req, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				return spec.Patch(service, args[0], &req, patchOpts.UpdateMask)
			})
		},
	}
	addPatchMutationFlags(patch, &patchOpts)

	var archiveApply bool
	archive := &cobra.Command{
		Use:   spec.DeleteUse,
		Short: "Archive the resource",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMutation(ctx, archiveApply, spec.DeleteOperation, args[0], nil, func() (any, error) {
				service, err := ctx.AdminService(cmd.Context())
				if err != nil {
					return nil, err
				}
				if spec.ArchiveDimension != nil {
					return spec.ArchiveDimension(service, args[0])
				}
				return spec.ArchiveMetric(service, args[0])
			})
		},
	}
	addApplyFlag(archive, &archiveApply)
	parent.AddCommand(create, patch, archive)
}
