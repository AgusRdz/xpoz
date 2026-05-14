#!/bin/sh
# Sign an artifact with the Ed25519 private key.
# Usage: scripts/sign.sh <artifact> <signature-output>
# Requires: SIGNING_KEY env var (base64-encoded PEM private key)
set -e

ARTIFACT="$1"
SIGNATURE="$2"

if [ -z "$SIGNING_KEY" ]; then
  echo "WARNING: SIGNING_KEY not set — skipping signature for $ARTIFACT" >&2
  touch "$SIGNATURE"
  exit 0
fi

KEYFILE="/tmp/xpoz-signing-$$.pem"
echo "$SIGNING_KEY" | base64 -d > "$KEYFILE"
openssl pkeyutl -sign -inkey "$KEYFILE" -rawin -in "$ARTIFACT" \
  | xxd -p -c 256 | tr -d '\n' > "$SIGNATURE"
rm -f "$KEYFILE"
