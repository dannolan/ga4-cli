# GA4 CLI

Tiny, shippable Google Analytics 4 from your terminal.

`ga4` is a Go rewrite of the local Node GA4 wrapper. It is built for agents, scripts, and non-technical teammates who should not need Node, npm, TypeScript, or a dependency tree just to pull analytics.

## Why

- Single compiled binary.
- Reuses local Google OAuth credentials.
- JSON-native command surface for agents.
- Friendly table/CSV/Markdown output for humans.
- Mutation-safe by default. Mutating API calls dry-run until you pass `--apply`.

## Install

Homebrew:

```bash
brew tap dannolan/tap
brew install ga4-cli
ga4 --help
```

From source:

```bash
go build -o ga4 ./cmd/ga4
./ga4 --help
```

## Configure

The CLI reads credentials from environment variables, `~/.config/ga4-cli/env`, or `~/.config/ga4-cli/config.json`.

```bash
export CLIENT_ID="..."
export CLIENT_SECRET="..."
export GA4_PROPERTY_ID="526319832"
```

It reuses OAuth tokens from:

- `~/.config/ga4-cli/token.json`
- legacy `~/.ga4-cli/token.json`

First-time setup:

```bash
ga4 config init \
  --client-id "$CLIENT_ID" \
  --client-secret "$CLIENT_SECRET" \
  --property "$GA4_PROPERTY_ID"

ga4 auth login
ga4 smoke
```

`ga4 auth login` asks Google for the scopes used by the official GA4 Data and Admin clients:

- `https://www.googleapis.com/auth/analytics.readonly`
- `https://www.googleapis.com/auth/analytics.edit`

The CLI can wrap both read-only and mutating Admin operations. Mutations are dry-run by default and require `--apply`; Google still requires the broader Admin scope for some read methods such as change history.

## Use

Quick checks:

```bash
ga4 doctor
ga4 smoke
ga4 events --limit 10
ga4 pages --format csv
ga4 report -m users,sessions -d country --format json
```

Raw Data API calls:

```bash
ga4 data metadata
ga4 data realtime -m activeUsers
ga4 data compatibility -m sessions -d country --compatibility COMPATIBLE

printf '%s' '{"dateRanges":[{"startDate":"2026-05-20","endDate":"2026-05-21"}],"metrics":[{"name":"sessions"}],"limit":"1"}' \
  | ga4 data run-report
```

Admin API read-only calls:

```bash
ga4 admin account-summaries list
ga4 admin accounts list
ga4 admin properties get properties/123456789
ga4 admin property-resources data-streams list properties/123456789
```

Admin commands require `analyticsadmin.googleapis.com` to be enabled on the OAuth project.
Run `ga4 doctor` if Admin commands fail; it reports the exact Google enablement URL when the service is disabled.

Mutating calls use the same API-shaped JSON bodies and dry-run by default:

```bash
printf '%s' '{"displayName":"New Property"}' \
  | ga4 admin properties create

printf '%s' '{"displayName":"New Property"}' \
  | ga4 admin properties create --apply
```

Use `ga4 manifest` to list which commands are mutations.

## Output

Report commands support:

- `--format table`
- `--format json`
- `--format markdown`
- `--format csv`

Metadata and raw API commands emit JSON.

## Coverage

See [docs/api-coverage.md](docs/api-coverage.md).

## Verify

```bash
go test ./...
go build -o ga4 ./cmd/ga4
scripts/smoke-readonly.sh
```
