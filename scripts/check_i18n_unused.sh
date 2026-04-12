#!/usr/bin/env bash
# check_i18n_unused.sh
#
# Finds keys in i18n/zh-CN.json that are not referenced anywhere in the Go
# source tree.
#
# A key is considered "referenced" if its exact string value appears in any
# .go file under internal/ — either as a Slug* constant value in slugs.go
# (compile-time constant) or as a raw string literal in older/generated code.
#
# Usage (run from project root):
#   ./scripts/check_i18n_unused.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

I18N_FILE="$ROOT/i18n/zh-CN.json"
SRC_DIR="$ROOT/internal"

if [[ ! -f "$I18N_FILE" ]]; then
    echo "ERROR: $I18N_FILE not found" >&2
    exit 1
fi

if ! command -v jq &>/dev/null; then
    echo "ERROR: jq is required but not installed (brew install jq)" >&2
    exit 1
fi

total=0
unused_count=0
unused_keys=""

while IFS= read -r key; do
    total=$((total + 1))
    # Search all .go files for the exact key string in double-quotes.
    if ! grep -rq "\"${key}\"" "$SRC_DIR" --include="*.go" 2>/dev/null; then
        unused_count=$((unused_count + 1))
        unused_keys="${unused_keys}  - ${key}\n"
    fi
done < <(jq -r 'keys[]' "$I18N_FILE")

if [[ $unused_count -eq 0 ]]; then
    echo "✓ All $total i18n keys are referenced in Go source."
    exit 0
fi

echo "✗ Found $unused_count unreferenced i18n key(s) in i18n/zh-CN.json:"
echo ""
printf "%b" "$unused_keys"
echo ""
echo "These keys have no corresponding Slug* constant or inline string in internal/."
echo "Consider removing them from i18n/zh-CN.json or adding the missing code references."
exit 1
