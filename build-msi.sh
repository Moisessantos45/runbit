#!/usr/bin/env bash
# =============================================================
#  build-msi.sh — Compila runbit para Windows y genera el .msi
#  Uso: ./build-msi.sh [version]
#  Ejemplo: ./build-msi.sh 1.2.0
# =============================================================
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$PROJECT_ROOT/dist"

# Versión: argumento CLI o fallback a "1.0.0"
VERSION="${1:-1.0.0}"
MSI_NAME="runbit_${VERSION}_amd64.msi"

echo "==> [1/3] Compilando runbit v${VERSION} para Windows con Wails..."
cd "$PROJECT_ROOT"
# -skipbindings para acelerar si ya se generaron o no son necesarios
wails build -clean \
  -platform windows/amd64 \
  -tags webkit2_41 \
  -ldflags "-X main.Version=${VERSION}" \
  -upx \
  -upxflags "--best --lzma"

echo "==> Binario generado en: build/bin/runbit.exe"

echo "==> [2/3] Generando archivo WXS para wixl..."
mkdir -p "$PROJECT_ROOT/packaging"
cat <<EOF > "$PROJECT_ROOT/packaging/runbit.wxs"
<?xml version="1.0" encoding="utf-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="Runbit" Language="1033" Version="${VERSION}" Manufacturer="Moisessantos45" UpgradeCode="12345678-1234-1234-1234-123456789012">
    <Package Description="Runbit Installer" InstallerVersion="200" Compressed="yes" InstallScope="perMachine" />
    <Media Id="1" Cabinet="runbit.cab" EmbedCab="yes" />

    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFilesFolder">
        <Directory Id="INSTALLDIR" Name="Runbit">
          <Component Id="RunbitExecutable" Guid="*">
            <File Id="RunbitExe" Source="$PROJECT_ROOT/build/bin/runbit.exe" KeyPath="yes" />
          </Component>
        </Directory>
      </Directory>
    </Directory>

    <Feature Id="MainApplication" Title="Main Application" Level="1">
      <ComponentRef Id="RunbitExecutable" />
    </Feature>
  </Product>
</Wix>
EOF

echo "==> [3/3] Generando el paquete .msi (v${VERSION}) con wixl..."
mkdir -p "$DIST_DIR"

cd "$PROJECT_ROOT/packaging"
if command -v wixl &>/dev/null; then
  wixl -o "$DIST_DIR/$MSI_NAME" runbit.wxs
  echo ""
  echo "Paquete MSI listo: dist/$MSI_NAME"
else
  echo "ERROR: wixl no está instalado. Asegúrate de tener wixl instalado:"
  echo "  sudo apt install wixl"
  exit 1
fi
