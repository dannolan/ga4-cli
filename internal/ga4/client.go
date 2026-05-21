package ga4

import (
	"context"
	"fmt"
	"strings"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

type Client struct {
	service *analyticsdata.Service
}

func NewClient(ctx context.Context, httpClient option.ClientOption) (*Client, error) {
	service, err := analyticsdata.NewService(ctx, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{service: service}, nil
}

func (c *Client) RunReport(ctx context.Context, params ReportParams) (ReportResult, error) {
	if strings.TrimSpace(params.PropertyID) == "" {
		return ReportResult{}, fmt.Errorf("property ID is required")
	}
	if len(params.Metrics) == 0 {
		return ReportResult{}, fmt.Errorf("at least one metric is required")
	}
	req := &analyticsdata.RunReportRequest{
		DateRanges: toAPIDateRanges(params.DateRanges),
		Dimensions: toAPIDimensions(params.Dimensions),
		Metrics:    toAPIMetrics(params.Metrics),
		Limit:      params.Limit,
		Offset:     params.Offset,
		OrderBys:   toAPIOrderBys(params.OrderBys),
	}
	if req.Limit == 0 {
		req.Limit = 100
	}
	if len(req.DateRanges) == 0 {
		return ReportResult{}, fmt.Errorf("at least one date range is required")
	}
	property := "properties/" + strings.TrimPrefix(params.PropertyID, "properties/")
	response, err := c.service.Properties.RunReport(property, req).Context(ctx).Do()
	if err != nil {
		return ReportResult{}, fmt.Errorf("run report: %w", err)
	}
	return fromAPIReport(response), nil
}

func (c *Client) TopPages(ctx context.Context, propertyID string, dateRange DateRange, limit int64) (ReportResult, error) {
	return c.RunReport(ctx, ReportParams{
		PropertyID: propertyID,
		Dimensions: []string{"pagePath", "pageTitle"},
		Metrics:    []string{"screenPageViews", "sessions"},
		DateRanges: []DateRange{dateRange},
		Limit:      limit,
	})
}

func (c *Client) TopEvents(ctx context.Context, propertyID string, dateRange DateRange, limit int64) (ReportResult, error) {
	return c.RunReport(ctx, ReportParams{
		PropertyID: propertyID,
		Dimensions: []string{"eventName"},
		Metrics:    []string{"eventCount"},
		DateRanges: []DateRange{dateRange},
		Limit:      limit,
	})
}

func (c *Client) RawRunReport(ctx context.Context, propertyID string, req *analyticsdata.RunReportRequest) (*analyticsdata.RunReportResponse, error) {
	return c.service.Properties.RunReport(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) BatchRunReports(ctx context.Context, propertyID string, req *analyticsdata.BatchRunReportsRequest) (*analyticsdata.BatchRunReportsResponse, error) {
	return c.service.Properties.BatchRunReports(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) RunPivotReport(ctx context.Context, propertyID string, req *analyticsdata.RunPivotReportRequest) (*analyticsdata.RunPivotReportResponse, error) {
	return c.service.Properties.RunPivotReport(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) BatchRunPivotReports(ctx context.Context, propertyID string, req *analyticsdata.BatchRunPivotReportsRequest) (*analyticsdata.BatchRunPivotReportsResponse, error) {
	return c.service.Properties.BatchRunPivotReports(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) RunRealtimeReport(ctx context.Context, propertyID string, req *analyticsdata.RunRealtimeReportRequest) (*analyticsdata.RunRealtimeReportResponse, error) {
	return c.service.Properties.RunRealtimeReport(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) CheckCompatibility(ctx context.Context, propertyID string, req *analyticsdata.CheckCompatibilityRequest) (*analyticsdata.CheckCompatibilityResponse, error) {
	return c.service.Properties.CheckCompatibility(propertyResource(propertyID), req).Context(ctx).Do()
}

func (c *Client) Metadata(ctx context.Context, propertyID string) (*analyticsdata.Metadata, error) {
	return c.service.Properties.GetMetadata(propertyResource(propertyID) + "/metadata").Context(ctx).Do()
}

func (c *Client) AudienceExportsList(ctx context.Context, propertyID string, pageSize int64) (*analyticsdata.ListAudienceExportsResponse, error) {
	call := c.service.Properties.AudienceExports.List(propertyResource(propertyID)).Context(ctx)
	if pageSize > 0 {
		call.PageSize(pageSize)
	}
	return call.Do()
}

func (c *Client) AudienceExportGet(ctx context.Context, name string) (*analyticsdata.AudienceExport, error) {
	return c.service.Properties.AudienceExports.Get(name).Context(ctx).Do()
}

func (c *Client) AudienceExportQuery(ctx context.Context, name string, req *analyticsdata.QueryAudienceExportRequest) (*analyticsdata.QueryAudienceExportResponse, error) {
	return c.service.Properties.AudienceExports.Query(name, req).Context(ctx).Do()
}

func propertyResource(propertyID string) string {
	return "properties/" + strings.TrimPrefix(propertyID, "properties/")
}

func toAPIDateRanges(ranges []DateRange) []*analyticsdata.DateRange {
	out := make([]*analyticsdata.DateRange, 0, len(ranges))
	for _, r := range ranges {
		out = append(out, &analyticsdata.DateRange{StartDate: r.StartDate, EndDate: r.EndDate})
	}
	return out
}

func toAPIDimensions(values []string) []*analyticsdata.Dimension {
	out := make([]*analyticsdata.Dimension, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, &analyticsdata.Dimension{Name: strings.TrimSpace(value)})
		}
	}
	return out
}

func toAPIMetrics(values []string) []*analyticsdata.Metric {
	out := make([]*analyticsdata.Metric, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, &analyticsdata.Metric{Name: strings.TrimSpace(value)})
		}
	}
	return out
}

func toAPIOrderBys(values []OrderBy) []*analyticsdata.OrderBy {
	out := make([]*analyticsdata.OrderBy, 0, len(values))
	for _, value := range values {
		field := strings.TrimSpace(value.Field)
		if field == "" {
			continue
		}
		out = append(out, &analyticsdata.OrderBy{
			Desc:   value.Desc,
			Metric: &analyticsdata.MetricOrderBy{MetricName: field},
		})
	}
	return out
}

func fromAPIReport(response *analyticsdata.RunReportResponse) ReportResult {
	out := ReportResult{
		RowCount:         response.RowCount,
		Rows:             []ReportRow{},
		DimensionHeaders: []string{},
		MetricHeaders:    []string{},
	}
	for _, header := range response.DimensionHeaders {
		out.DimensionHeaders = append(out.DimensionHeaders, header.Name)
	}
	for _, header := range response.MetricHeaders {
		out.MetricHeaders = append(out.MetricHeaders, header.Name)
	}
	for _, row := range response.Rows {
		item := ReportRow{}
		for _, value := range row.DimensionValues {
			item.Dimensions = append(item.Dimensions, value.Value)
		}
		for _, value := range row.MetricValues {
			item.Metrics = append(item.Metrics, value.Value)
		}
		out.Rows = append(out.Rows, item)
	}
	return out
}
