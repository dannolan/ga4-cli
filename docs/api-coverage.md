# API Coverage

This CLI targets Google Analytics Data API `v1beta` and Google Analytics Admin API `v1beta` using the official generated Go clients.

Mutating commands are implemented but dry-run by default. They print the request they would send and require `--apply` before calling Google.

## Data API Read-Only Coverage

Implemented:

- `properties.getMetadata` via `ga4 data metadata`
- `properties.runReport` via `ga4 report` and `ga4 data run-report`
- `properties.batchRunReports` via `ga4 data batch-run-reports`
- `properties.runPivotReport` via `ga4 data run-pivot-report`
- `properties.batchRunPivotReports` via `ga4 data batch-run-pivot-reports`
- `properties.runRealtimeReport` via `ga4 data realtime`
- `properties.checkCompatibility` via `ga4 data compatibility`
- `properties.audienceExports.list` via `ga4 data audience-exports list`
- `properties.audienceExports.get` via `ga4 data audience-exports get`
- `properties.audienceExports.query` via `ga4 data audience-exports query`

Implemented as dry-run/apply mutation:

- `properties.audienceExports.create` via `ga4 data audience-exports create`

## Admin API Read-Only Coverage

Implemented:

- `accountSummaries.list`
- `accounts.list`
- `accounts.get`
- `accounts.getDataSharingSettings`
- `accounts.runAccessReport`
- `accounts.searchChangeHistoryEvents`
- `properties.list`
- `properties.get`
- `properties.getDataRetentionSettings`
- `properties.runAccessReport`
- `properties.conversionEvents.list/get`
- `properties.customDimensions.list/get`
- `properties.customMetrics.list/get`
- `properties.dataStreams.list/get`
- `properties.dataStreams.measurementProtocolSecrets.list/get`
- `properties.firebaseLinks.list`
- `properties.googleAdsLinks.list`
- `properties.keyEvents.list/get`

Implemented as dry-run/apply mutations:

- `accounts.delete`
- `accounts.patch`
- `accounts.provisionAccountTicket`
- `properties.acknowledgeUserDataCollection`
- `properties.create`
- `properties.delete`
- `properties.patch`
- `properties.updateDataRetentionSettings`
- `properties.conversionEvents.create/delete/patch`
- `properties.customDimensions.archive/create/patch`
- `properties.customMetrics.archive/create/patch`
- `properties.dataStreams.create/delete/patch`
- `properties.dataStreams.measurementProtocolSecrets.create/delete/patch`
- `properties.firebaseLinks.create/delete`
- `properties.googleAdsLinks.create/delete/patch`
- `properties.keyEvents.create/delete/patch`

## Verification

Run:

```bash
go test ./...
go build -o ga4 ./cmd/ga4
scripts/smoke-readonly.sh
```

The smoke script runs every Data API read-only method. Admin API checks run when `analyticsadmin.googleapis.com` is enabled for the OAuth project; otherwise the script records that Admin verification is skipped because the service is disabled.
