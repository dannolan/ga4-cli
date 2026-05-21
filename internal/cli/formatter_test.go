package cli

import (
	"encoding/json"
	"testing"

	"github.com/dannolan/ga4-cli/internal/ga4"
)

func TestFormatJSONUsesEmptyHeaderArrays(t *testing.T) {
	text, err := formatReport(ga4.ReportResult{
		Rows:          []ga4.ReportRow{{Metrics: []string{"25"}}},
		MetricHeaders: []string{"sessions"},
		RowCount:      1,
	}, "json")
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatal(err)
	}
	summary := payload["summary"].(map[string]any)
	if summary["dimensionHeaders"] == nil {
		t.Fatalf("dimensionHeaders was nil: %s", text)
	}
}

func TestFormatCSVQuotesCommas(t *testing.T) {
	text, err := formatReport(ga4.ReportResult{
		Rows:             []ga4.ReportRow{{Dimensions: []string{"/,x"}, Metrics: []string{"1"}}},
		DimensionHeaders: []string{"pagePath"},
		MetricHeaders:    []string{"sessions"},
		RowCount:         1,
	}, "csv")
	if err != nil {
		t.Fatal(err)
	}
	want := "pagePath,sessions\n\"/,x\",1"
	if text != want {
		t.Fatalf("got %q, want %q", text, want)
	}
}
