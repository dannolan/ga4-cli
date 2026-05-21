# Verification Record

Last verified: 2026-05-21

## Local Build

```bash
go test ./...
go build -o ga4 ./cmd/ga4
```

Result: passed.

## Read-Only Smoke Suite

```bash
scripts/smoke-readonly.sh
```

Result: passed for all GA4 Data API read-only methods.

Verified live with local credentials:

- `data metadata`
- `data realtime`
- `data compatibility`
- `data run-report`
- `data batch-run-reports`
- `data run-pivot-report`
- `data batch-run-pivot-reports`
- `data audience-exports list`

Audience export `get` and `query` are implemented but skipped by the smoke script when the property has no existing audience exports.

## Admin API Status

Admin API commands are implemented, but live verification is blocked by the OAuth project configuration. Google returns:

```text
Google Analytics Admin API has not been used in project 1013065713343 before or it is disabled.
service: analyticsadmin.googleapis.com
reason: SERVICE_DISABLED
```

Enable URL:

```text
https://console.developers.google.com/apis/api/analyticsadmin.googleapis.com/overview?project=1013065713343
```

After enabling the service, rerun:

```bash
scripts/smoke-readonly.sh
```

The script will automatically run the Admin API checks instead of skipping them.

## Homebrew

```bash
brew audit --strict --online dannolan/tap/ga4-cli
brew test dannolan/tap/ga4-cli
brew install dannolan/tap/ga4-cli
ga4 smoke
```

Result: passed.

Installed binary:

```json
{
  "version": "0.1.0",
  "commit": "25e28e4"
}
```
