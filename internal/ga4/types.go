package ga4

type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ReportParams struct {
	PropertyID      string      `json:"propertyId"`
	Dimensions      []string    `json:"dimensions,omitempty"`
	Metrics         []string    `json:"metrics"`
	DateRanges      []DateRange `json:"dateRanges"`
	Limit           int64       `json:"limit,omitempty"`
	Offset          int64       `json:"offset,omitempty"`
	OrderBys        []OrderBy   `json:"orderBys,omitempty"`
	DimensionFilter string      `json:"dimensionFilter,omitempty"`
	MetricFilter    string      `json:"metricFilter,omitempty"`
}

type OrderBy struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc,omitempty"`
}

type ReportRow struct {
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}

type ReportResult struct {
	Rows             []ReportRow `json:"rows"`
	DimensionHeaders []string    `json:"dimensionHeaders"`
	MetricHeaders    []string    `json:"metricHeaders"`
	RowCount         int64       `json:"rowCount"`
}
