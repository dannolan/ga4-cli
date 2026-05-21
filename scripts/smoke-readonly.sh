#!/usr/bin/env bash
set -euo pipefail

bin="${BIN:-./ga4}"
property="${GA4_PROPERTY_ID:-}"

if [[ -z "$property" ]]; then
  property="$("$bin" config show --pretty=false | jq -r '.property_id')"
fi

if [[ -z "$property" || "$property" == "null" ]]; then
  echo "missing GA4_PROPERTY_ID" >&2
  exit 1
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

echo "data: metadata"
"$bin" data metadata --property "$property" --pretty=false | jq -e '.dimensions and .metrics' >/dev/null

echo "data: realtime"
"$bin" data realtime --property "$property" -m activeUsers --pretty=false | jq -e '.kind == "analyticsData#runRealtimeReport"' >/dev/null

echo "data: compatibility"
"$bin" data compatibility --property "$property" -m sessions -d country --compatibility COMPATIBLE --pretty=false | jq -e '.dimensionCompatibilities and .metricCompatibilities' >/dev/null

cat >"$tmpdir/report.json" <<JSON
{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"metrics":[{"name":"sessions"}],"limit":"1"}
JSON

echo "data: run-report"
"$bin" data run-report --property "$property" --body "$tmpdir/report.json" --pretty=false | jq -e '.kind == "analyticsData#runReport"' >/dev/null

cat >"$tmpdir/batch.json" <<JSON
{"requests":[{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"metrics":[{"name":"sessions"}],"limit":"1"}]}
JSON

echo "data: batch-run-reports"
"$bin" data batch-run-reports --property "$property" --body "$tmpdir/batch.json" --pretty=false | jq -e '.kind == "analyticsData#batchRunReports"' >/dev/null

cat >"$tmpdir/pivot.json" <<JSON
{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"dimensions":[{"name":"country"}],"metrics":[{"name":"sessions"}],"pivots":[{"fieldNames":["country"],"limit":"5"}]}
JSON

echo "data: run-pivot-report"
"$bin" data run-pivot-report --property "$property" --body "$tmpdir/pivot.json" --pretty=false | jq -e '.kind == "analyticsData#runPivotReport"' >/dev/null

cat >"$tmpdir/batch-pivot.json" <<JSON
{"requests":[{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"dimensions":[{"name":"country"}],"metrics":[{"name":"sessions"}],"pivots":[{"fieldNames":["country"],"limit":"5"}]}]}
JSON

echo "data: batch-run-pivot-reports"
"$bin" data batch-run-pivot-reports --property "$property" --body "$tmpdir/batch-pivot.json" --pretty=false | jq -e '.kind == "analyticsData#batchRunPivotReports"' >/dev/null

echo "data: audience-exports list"
audience_exports="$("$bin" data audience-exports list --property "$property" --pretty=false)"
echo "$audience_exports" | jq -e 'type == "object"' >/dev/null
first_export="$(echo "$audience_exports" | jq -r '.audienceExports[0].name // empty')"
if [[ -n "$first_export" ]]; then
  echo "data: audience-exports get/query"
  "$bin" data audience-exports get "$first_export" --pretty=false | jq -e '.name' >/dev/null
  "$bin" data audience-exports query "$first_export" --limit 1 --pretty=false | jq -e '.kind == "analyticsData#queryAudienceExport"' >/dev/null
else
  echo "data: audience-exports get/query skipped; no exports exist"
fi

echo "admin: account-summaries"
admin_output="$("$bin" admin account-summaries list --page-size 5 --pretty=false 2>&1)" || {
  if grep -q 'analyticsadmin.googleapis.com.*disabled\|SERVICE_DISABLED' <<<"$admin_output"; then
    echo "admin: skipped; Analytics Admin API is disabled for this OAuth project"
    echo "readonly smoke passed for Data API"
    exit 0
  fi
  echo "$admin_output" >&2
  exit 1
}
echo "$admin_output" | jq -e '.accountSummaries' >/dev/null

account="$(echo "$admin_output" | jq -r '.accountSummaries[0].account // empty')"
summary_property="$(echo "$admin_output" | jq -r '.accountSummaries[0].propertySummaries[0].property // empty')"
property_resource="properties/$property"
if [[ -n "$summary_property" ]]; then
  property_resource="$summary_property"
fi

if [[ -n "$account" ]]; then
  echo "admin: accounts list/get/settings/change-history/access-report"
  "$bin" admin accounts list --page-size 5 --pretty=false | jq -e '.accounts' >/dev/null
  "$bin" admin accounts get "$account" --pretty=false | jq -e '.name' >/dev/null
  "$bin" admin accounts data-sharing-settings "$account" --pretty=false | jq -e '.name' >/dev/null
  printf '%s' '{"pageSize":1}' | "$bin" admin accounts change-history "$account" --pretty=false | jq -e '.changeHistoryEvents != null or .kind' >/dev/null
  printf '%s' '{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"metrics":[{"metricName":"accessCount"}],"limit":"1"}' | "$bin" admin accounts access-report "$account" --pretty=false | jq -e '.kind' >/dev/null
fi

echo "admin: properties list/get/retention/access-report"
"$bin" admin properties list --filter "parent:$account" --page-size 5 --pretty=false | jq -e '.properties' >/dev/null
"$bin" admin properties get "$property_resource" --pretty=false | jq -e '.name' >/dev/null
"$bin" admin properties data-retention-settings "$property_resource" --pretty=false | jq -e '.name' >/dev/null
printf '%s' '{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"metrics":[{"metricName":"accessCount"}],"limit":"1"}' | "$bin" admin properties access-report "$property_resource" --pretty=false | jq -e '.kind' >/dev/null

for resource in conversion-events custom-dimensions custom-metrics data-streams firebase-links google-ads-links key-events; do
  echo "admin: property-resources $resource list"
  output="$("$bin" admin property-resources "$resource" list "$property_resource" --pretty=false)"
  echo "$output" | jq -e 'type == "object"' >/dev/null
done

data_stream="$("$bin" admin property-resources data-streams list "$property_resource" --pretty=false | jq -r '.dataStreams[0].name // empty')"
if [[ -n "$data_stream" ]]; then
  echo "admin: measurement-protocol-secrets list"
  "$bin" admin property-resources measurement-protocol-secrets list "$data_stream" --pretty=false | jq -e 'type == "object"' >/dev/null
fi

echo "readonly smoke passed"
