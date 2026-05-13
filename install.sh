#!/bin/sh
set -e

REPO="AgusRdz/xpoz"
BINARY="xpoz"

# ── platform detection ────────────────────────────────────────────────────────

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)          ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "✗ Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  linux | darwin) ;;
  *)
    echo "✗ Unsupported OS: $OS (use install.ps1 on Windows)" >&2
    exit 1
    ;;
esac

# ── install directory ─────────────────────────────────────────────────────────

if [ -n "$XPOZ_INSTALL_DIR" ]; then
  INSTALL_DIR="$XPOZ_INSTALL_DIR"
elif [ -w /usr/local/bin ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
fi

mkdir -p "$INSTALL_DIR"

# ── latest version ────────────────────────────────────────────────────────────

VERSION="${XPOZ_VERSION:-}"
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | sed 's/.*"tag_name": *"//;s/".*//')"
fi

if [ -z "$VERSION" ]; then
  echo "✗ Could not determine latest version. Set XPOZ_VERSION to override." >&2
  exit 1
fi

# ── download and extract ──────────────────────────────────────────────────────

ARCHIVE="${BINARY}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
DEST="${INSTALL_DIR}/${BINARY}"
TMPDIR="$(mktemp -d)"

echo "  Downloading xpoz ${VERSION} (${OS}/${ARCH})..."
curl -fsSL "$URL" | tar -xz -C "$TMPDIR"
mv "$TMPDIR/$BINARY" "$DEST"
chmod +x "$DEST"
rm -rf "$TMPDIR"

# ── PATH setup ────────────────────────────────────────────────────────────────

case ":$PATH:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    SHELL_NAME="$(basename "${SHELL:-sh}")"
    RC_FILE=""
    case "$SHELL_NAME" in
      zsh)  RC_FILE="$HOME/.zshrc" ;;
      bash) RC_FILE="$HOME/.bashrc" ;;
    esac

    if [ -n "$RC_FILE" ]; then
      EXPORT_LINE="export PATH=\"${INSTALL_DIR}:\$PATH\""
      if ! grep -qF "$EXPORT_LINE" "$RC_FILE" 2>/dev/null; then
        printf '\n# xpoz\n%s\n' "$EXPORT_LINE" >> "$RC_FILE"
        echo "  Added ${INSTALL_DIR} to PATH in ${RC_FILE}"
      fi
    fi

    # Make available in the current session without a shell reload
    export PATH="${INSTALL_DIR}:${PATH}"
    ;;
esac

echo "✓ xpoz ${VERSION} installed to ${DEST}"
echo ""
echo "  Run 'xpoz setup' to get started."
