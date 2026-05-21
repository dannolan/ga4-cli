#!/usr/bin/env bash
set -euo pipefail

version="${1:?usage: scripts/release.sh vX.Y.Z}"
commit="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
date="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
dist="dist"

rm -rf "$dist"
mkdir -p "$dist"

ldflags="-s -w -X github.com/dannolan/ga4-cli/internal/cli.Version=${version#v} -X github.com/dannolan/ga4-cli/internal/cli.Commit=${commit} -X github.com/dannolan/ga4-cli/internal/cli.Date=${date}"

for target in darwin/arm64 darwin/amd64 linux/arm64 linux/amd64; do
  goos="${target%/*}"
  goarch="${target#*/}"
  work="$dist/ga4_${goos}_${goarch}"
  mkdir -p "$work"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -ldflags "$ldflags" -o "$work/ga4" ./cmd/ga4
  (cd "$work" && zip -q "../ga4_${goos}_${goarch}.zip" ga4)
done

(cd "$dist" && shasum -a 256 *.zip | tee SHA256SUMS)
