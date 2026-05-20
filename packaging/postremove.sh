#!/bin/sh
set -e
update-desktop-database /usr/share/applications || true
gtk-update-icon-cache /usr/share/icons/hicolor || true
