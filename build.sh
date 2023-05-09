#!/bin/bash

set -euo pipefail

set -x

TARGET=ios,iossimulator,macos
OUTPUT=IEnvoyProxy.xcframework
TEMPDIR="${TEMPDIR:/tmp/}"
TEMPDIR="${TMPDIR}IEnvoyProxy"

echo "TEMPDIR: ${TEMPDIR}"

if test "$1" = "android"; then
  TARGET=android
  OUTPUT=IEnvoyProxy.aar
fi

cd "$(dirname "$0")" || exit 1

if test -e $OUTPUT; then
    echo "--- No build necessary, $OUTPUT already exists."
    exit
fi

# Install dependencies. Go itself is a prerequisite.
printf '\n--- Golang 1.16 or up needs to be installed! Try "brew install go" on MacOS or "snap install go" on Linux if we fail further down!'
printf '\n--- Installing gomobile...\n'
go install golang.org/x/mobile/cmd/gomobile@latest

# Prepare build environment
# Go leaks the build path in to the binary, so use a temp dir to build
# based on https://github.com/tladesignz/IPtProxy/pull/38
printf '\n\n--- Prepare build environment at %s...\n' "$TEMPDIR"
CURRENT=$PWD
rm -rf "$TEMPDIR" || true
mkdir -p "$TEMPDIR"
cp -a IEnvoyProxy "$TEMPDIR/"

# Fetch submodules.
printf '\n\n--- Fetching submodule dependencies...\n'
if test -e ".git"; then
    # There's a .git directory - we must be in the development pod.
    git submodule update --init --recursive
    cd hysteria || exit 1
    git reset --hard
    cp -a . "$TEMPDIR/hysteria"
    cd ../v2ray-core || exit 1
    git reset --hard
    git clean -ffd # we add a file
    cp -a . "$TEMPDIR/v2ray-core"
    cd ../snowflake || exit 1
    git reset --hard
    cp -a . "$TEMPDIR/snowflake"
    cd ..
else
    # No .git directory - That's a normal install.
    git clone https://github.com/HyNetwork/hysteria.git "$TEMPDIR/hysteria"
    cd hysteria || exit 1
    git checkout --force --quiet da16c88
    cd ..
    git clone https://github.com/v2fly/v2ray-core.git "$TEMPDIR/v2ray-core"
    cd v2ray-core || exit 1
    git checkout --force --quiet b4069f74
    cd ..
    git clone https://git.torproject.org/pluggable-transports/snowflake.git "$TEMPDIR/snowflake"
    cd "$TEMPDIR/snowflake" || exit 1
    git checkout --force --quiet 7b77001
    cd "$CURRENT" || exit 1
fi

# Apply patches.
printf '\n\n--- Apply patches to submodules...\n'
echo `pwd`
patch --directory=$TEMPDIR/hysteria --strip=1 < hysteria.patch
patch --directory=$TEMPDIR/v2ray-core --strip=1 < v2ray-core.patch
patch --directory=$TEMPDIR/snowflake --strip=1 < snowflake.patch

# Compile framework.
printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
cd "${TEMPDIR}/IEnvoyProxy" || exit 1

gomobile init

MACOSX_DEPLOYMENT_TARGET=11.0 gomobile bind -target=$TARGET -o $CURRENT/$OUTPUT -iosversion=11.0 -androidapi=19 -v -tags=netcgo -trimpath

rm -rf "$TEMPDIR"

printf '\n\n--- Done.\n\n'
