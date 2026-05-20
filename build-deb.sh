#!/usr/bin/env bash
# =============================================================
#  build-deb.sh — Compila runbit y genera el .deb
#  Uso: ./build-deb.sh [version]
#  Ejemplo: ./build-deb.sh 1.2.0
# =============================================================
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$PROJECT_ROOT/dist"

# Versión: argumento CLI o fallback a "1.0.0"
VERSION="${1:-1.0.0}"
DEB_NAME="runbit_${VERSION}_amd64.deb"

echo "==> [1/3] Compilando runbit v${VERSION} con Wails (webkit2_41 + UPX)..."
cd "$PROJECT_ROOT"
wails build -clean \
  -tags webkit2_41 \
  -ldflags "-X main.Version=${VERSION}" \
  -upx \
  -upxflags "--best --lzma"

echo "==> Binario generado en: build/bin/runbit"

echo "==> [2/3] Verificando nfpm..."
if ! command -v nfpm &>/dev/null; then
  echo "    nfpm no encontrado. Instalando..."
  go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
fi

echo "==> [3/3] Generando el paquete .deb (v${VERSION})..."
mkdir -p "$DIST_DIR"

chmod +x "$PROJECT_ROOT/packaging/postinstall.sh"
chmod +x "$PROJECT_ROOT/packaging/postremove.sh"

NFPM_TMP="$PROJECT_ROOT/packaging/nfpm_build.yaml"
sed "s/^version:.*/version: \"${VERSION}\"/" \
  "$PROJECT_ROOT/packaging/nfpm.yaml" > "$NFPM_TMP"

cd "$PROJECT_ROOT/packaging"
nfpm pkg --config nfpm_build.yaml --packager deb --target "$DIST_DIR/$DEB_NAME"
rm -f "$NFPM_TMP"

echo ""
echo "Paquete listo: dist/$DEB_NAME"
echo ""
echo "Para instalar:"
echo "  sudo dpkg -i dist/$DEB_NAME"
