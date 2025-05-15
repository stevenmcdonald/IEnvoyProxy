#!/bin/sh

cd "$(dirname "$0")" || exit 1

cd hysteria || exit 1
git reset --hard
cd ../v2ray-core || exit 1
git reset --hard
git clean -fdx
cd ..

patch --directory="hysteria" --strip=1 < hysteria.patch
patch --directory="v2ray-core" --strip=1 < v2ray-core.patch

cd "IEnvoyProxy" || exit 1

gomobile init

MACOSX_DEPLOYMENT_TARGET=11.0 gomobile bind -target=ios,iossimulator,macos -o "../IEnvoyProxy.xcframework" -iosversion=12.0 -v -tags=netcgo -trimpath
