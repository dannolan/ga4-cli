package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dannolan/ga4-cli/internal/ga4"
)

func formatReport(result ga4.ReportResult, format string) (string, error) {
	switch strings.ToLower(format) {
	case "table":
		return formatTable(result), nil
	case "json":
		return formatJSON(result)
	case "markdown":
		return formatMarkdown(result), nil
	case "csv":
		return formatCSV(result)
	default:
		return "", fmt.Errorf("unsupported format %q; use table, json, markdown, or csv", format)
	}
}

func formatComparison(current, previous ga4.ReportResult, format string) (string, error) {
	if strings.ToLower(format) == "json" {
		currentData, err := reportJSONValue(current)
		if err != nil {
			return "", err
		}
		previousData, err := reportJSONValue(previous)
		if err != nil {
			return "", err
		}
		data, err := json.MarshalIndent(map[string]any{
			"current period":  currentData,
			"previous period": previousData,
		}, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	currentText, err := formatReport(current, format)
	if err != nil {
		return "", err
	}
	previousText, err := formatReport(previous, format)
	if err != nil {
		return "", err
	}
	if strings.ToLower(format) == "markdown" {
		return "## Current Period\n" + currentText + "\n\n## Previous Period\n" + previousText, nil
	}
	return "=== Current Period ===\n" + currentText + "\n\n=== Previous Period ===\n" + previousText, nil
}

func formatTable(result ga4.ReportResult) string {
	if len(result.Rows) == 0 {
		return "No data found"
	}
	headers := append(append([]string{}, result.DimensionHeaders...), result.MetricHeaders...)
	rows := reportRows(result)
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = max(10, len(header))
	}
	for _, row := range rows {
		for i, cell := range row {
			widths[i] = max(widths[i], len(cell))
		}
	}
	border := func(left, middle, right string) string {
		parts := make([]string, len(widths))
		for i, width := range widths {
			parts[i] = strings.Repeat("-", width+2)
		}
		return left + strings.Join(parts, middle) + right
	}
	line := []string{border("+", "+", "+")}
	headerCells := make([]string, len(headers))
	for i, header := range headers {
		headerCells[i] = " " + padRight(header, widths[i]) + " "
	}
	line = append(line, "|"+strings.Join(headerCells, "|")+"|", border("+", "+", "+"))
	for _, row := range rows {
		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = " " + padRight(cell, widths[i]) + " "
		}
		line = append(line, "|"+strings.Join(cells, "|")+"|")
	}
	line = append(line, border("+", "+", "+"))
	return strings.Join(line, "\n")
}

func formatJSON(result ga4.ReportResult) (string, error) {
	value, err := reportJSONValue(result)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func reportJSONValue(result ga4.ReportResult) (map[string]any, error) {
	data := make([]map[string]string, 0, len(result.Rows))
	headers := append(append([]string{}, result.DimensionHeaders...), result.MetricHeaders...)
	for _, row := range reportRows(result) {
		item := map[string]string{}
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			}
		}
		data = append(data, item)
	}
	return map[string]any{
		"summary": map[string]any{
			"totalRows":        result.RowCount,
			"dimensionHeaders": nonNilStrings(result.DimensionHeaders),
			"metricHeaders":    nonNilStrings(result.MetricHeaders),
		},
		"data": data,
	}, nil
}

func formatMarkdown(result ga4.ReportResult) string {
	if len(result.Rows) == 0 {
		return "No data found"
	}
	headers := append(append([]string{}, result.DimensionHeaders...), result.MetricHeaders...)
	lines := []string{"# GA4 Report Results", "", "Total rows: " + strconv.FormatInt(result.RowCount, 10), ""}
	lines = append(lines, "| "+strings.Join(headers, " | ")+" |")
	separators := make([]string, len(headers))
	for i := range separators {
		separators[i] = "---"
	}
	lines = append(lines, "| "+strings.Join(separators, " | ")+" |")
	for _, row := range reportRows(result) {
		lines = append(lines, "| "+strings.Join(row, " | ")+" |")
	}
	return strings.Join(lines, "\n")
}

func formatCSV(result ga4.ReportResult) (string, error) {
	if len(result.Rows) == 0 {
		return "No data found", nil
	}
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	headers := append(append([]string{}, result.DimensionHeaders...), result.MetricHeaders...)
	if err := writer.Write(headers); err != nil {
		return "", err
	}
	for _, row := range reportRows(result) {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return strings.TrimRight(builder.String(), "\n"), writer.Error()
}

func writeOutput(text, path string) error {
	if path == "" {
		_, err := fmt.Fprintln(os.Stdout, text)
		return err
	}
	return os.WriteFile(path, []byte(text), 0o644)
}

func reportRows(result ga4.ReportResult) [][]string {
	rows := make([][]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		rows = append(rows, append(append([]string{}, row.Dimensions...), row.Metrics...))
	}
	return rows
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
