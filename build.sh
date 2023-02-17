#!/bin/sh

TARGET=ios,iossimulator,macos
OUTPUT=IEnvoyProxy.xcframework

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

# Fetch submodules.
printf '\n\n--- Fetching submodule dependencies...\n'
if test -e ".git"; then
    # There's a .git directory - we must be in the development pod.
    git submodule update --init --recursive
    cd dnstt || exit 1
    git reset --hard
    cd ../hysteria || exit 1
    git reset --hard
    cd ../v2ray-core || exit 1
    git reset --hard
    git clean -ffd
    cd ..
else
    # No .git directory - That's a normal install.
    git clone https://www.bamsoftware.com/git/dnstt.git
    cd dnstt || exit 1
    git checkout --force --quiet 04f04590
    cd ..
    git clone https://github.com/HyNetwork/hysteria.git
    cd hysteria || exit 1
    git checkout --force --quiet da16c88
    cd ..
    git clone https://github.com/v2fly/v2ray-core.git
    cd v2ray-core || exit 1
    git checkout --force --quiet b4069f74
    cd ..
fi

# Apply patches.
printf '\n\n--- Apply patches to submodules...\n'
echo `pwd`
patch --directory=dnstt --strip=1 < dnstt.patch
patch --directory=hysteria --strip=1 < hysteria.patch
patch --directory=v2ray-core --strip=1 < v2ray-core.patch

# Compile framework.
printf '\n\n--- Compile %s...\n' "$OUTPUT"
export PATH=~/go/bin:$PATH
cd IEnvoyProxy || exit 1

gomobile init

gomobile bind -target=$TARGET -o ../$OUTPUT -iosversion 11.0 -androidapi 19 -v

printf '\n\n--- Done.\n\n'
