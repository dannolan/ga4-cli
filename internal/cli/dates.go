package cli

import (
	"time"

	"github.com/dannolan/ga4-cli/internal/ga4"
)

func defaultDateRange(start, end string) ga4.DateRange {
	now := time.Now()
	if end == "" {
		end = now.Format(time.DateOnly)
	}
	if start == "" {
		start = now.AddDate(0, 0, -7).Format(time.DateOnly)
	}
	return ga4.DateRange{StartDate: start, EndDate: end}
}

func defaultComparisonRanges(currentStart, currentEnd, previousStart, previousEnd string) (ga4.DateRange, ga4.DateRange) {
	if currentStart != "" && currentEnd != "" && previousStart != "" && previousEnd != "" {
		return ga4.DateRange{StartDate: currentStart, EndDate: currentEnd}, ga4.DateRange{StartDate: previousStart, EndDate: previousEnd}
	}
	today := time.Now()
	currentEndDate := today
	currentStartDate := today.AddDate(0, 0, -6)
	previousEndDate := currentStartDate.AddDate(0, 0, -1)
	previousStartDate := previousEndDate.AddDate(0, 0, -6)
	return ga4.DateRange{
			StartDate: currentStartDate.Format(time.DateOnly),
			EndDate:   currentEndDate.Format(time.DateOnly),
		}, ga4.DateRange{
			StartDate: previousStartDate.Format(time.DateOnly),
			EndDate:   previousEndDate.Format(time.DateOnly),
		}
}
