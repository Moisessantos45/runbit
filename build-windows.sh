#!/usr/bin/env bash
# =============================================================
#  build-windows.sh — Compila runbit para Windows y genera el .exe installer
#  Uso: ./build-windows.sh [version]
#  Ejemplo: ./build-windows.sh 1.5.0
# =============================================================
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$PROJECT_ROOT/dist"

# Versión: argumento CLI o fallback a "1.0.0"
VERSION="${1:-1.0.0}"
SETUP_NAME="runbit_${VERSION}_amd64_setup.exe"

echo "==> [1/2] Compilando runbit v${VERSION} para Windows con Wails..."
cd "$PROJECT_ROOT"
wails build -clean \
  -platform windows/amd64 \
  -ldflags "-X main.Version=${VERSION}" \
  -upx \
  -upxflags "--best --lzma"

echo "==> Binario generado en: build/bin/runbit.exe"

echo "==> [2/2] Generando el instalador con NSIS (v${VERSION})..."
makensis -DAPPVERSION="${VERSION}" installer.nsi

echo "==> [3/3] Moviendo el instalador a la carpeta dist..."
mkdir -p "$DIST_DIR"
mv runbit-setup.exe "$DIST_DIR/$SETUP_NAME"

echo ""
echo "Instalador listo: dist/$SETUP_NAME"
