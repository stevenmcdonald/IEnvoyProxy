#!/bin/bash

set -eo pipefail

TARGET=ios,iossimulator,macos
OUTPUT=IEnvoyProxy.xcframework

# test if TMPDIR is unset: https://stackoverflow.com/a/13864829
if [[ -z ${TMPDIR} ]]; then
    # macOS
    TMPDIR=$(mktemp -dq)
else
    # Linux
    TMPDIR="${TMPDIR}IEnvoyProxy"
    mkdir ${TMPDIR} || true
fi
# TMPDIR may be unbound until now
set -u

# echo "TMPDIR: ${TMPDIR}"

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
printf '\n--- Golang 1.19 or up needs to be installed! Try "brew install go" on MacOS or "snap install go" on Linux if we fail further down!'
printf '\n--- Installing gomobile...\n'
go install golang.org/x/mobile/cmd/gomobile@latest

# Prepare build environment
# Go leaks the build path in to the binary, so use a temp dir to build
# based on https://github.com/tladesignz/IPtProxy/pull/38
printf '\n\n--- Prepare build environment at %s...\n' "$TMPDIR"
CURRENT=$PWD
rm -rf "$TMPDIR" || true
mkdir -p "$TMPDIR"
cp -a IEnvoyProxy "$TMPDIR/"

# Fetch submodules.
printf '\n\n--- Fetching submodule dependencies...\n'
if test -e ".git"; then
    # There's a .git directory - we must be in the development pod.
    git submodule update --init --recursive
    cd lyrebird || exit 1
    git reset --hard
    cp -a . "$TMPDIR/lyrebird"
    cd ../hysteria || exit 1
    git reset --hard
    cp -a . "$TMPDIR/hysteria"
    cd ../v2ray-core || exit 1
    git reset --hard
    git clean -ffd # we add a file
    cp -a . "$TMPDIR/v2ray-core"
    cd ../snowflake || exit 1
    git reset --hard
    cp -a . "$TMPDIR/snowflake"
    cd ..
else
    # No .git directory - That's a normal install.
    git clone https://gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/lyrebird.git "$TMPDIR/lyrebird"
    cd "$TMPDIR/lyrebird" || exit 1
    git checkout --force --quiet 3915dcd
    git clone https://github.com/HyNetwork/hysteria.git "$TMPDIR/hysteria"
    cd hysteria || exit 1
    git checkout --force --quiet b94f8a1
    cd ..
    git clone https://github.com/v2fly/v2ray-core.git "$TMPDIR/v2ray-core"
    cd v2ray-core || exit 1
    git checkout --force --quiet 9b526285
    cd ..
    git clone https://git.torproject.org/pluggable-transports/snowflake.git "$TMPDIR/snowflake"
    cd "$TMPDIR/snowflake" || exit 1
    git checkout --force --quiet 7b77001
    cd "$CURRENT" || exit 1
fi

# Apply patches.
printf '\n\n--- Apply patches to submodules...\n'
echo `pwd`
patch --directory=$TMPDIR/lyrebird --strip=1 < lyrebird.patch
patch --directory=$TMPDIR/hysteria --strip=1 < hysteria.patch
patch --directory=$TMPDIR/v2ray-core --strip=1 < v2ray-core.patch
patch --directory=$TMPDIR/snowflake --strip=1 < snowflake.patch

# Compile framework.
printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
cd "${TMPDIR}/IEnvoyProxy" || exit 1

gomobile init

MACOSX_DEPLOYMENT_TARGET=11.0 gomobile bind -target=$TARGET -o $CURRENT/$OUTPUT -iosversion=11.0 -androidapi=19 -v -tags=netcgo -trimpath

rm -rf "$TMPDIR"

printf '\n\n--- Done.\n\n'
