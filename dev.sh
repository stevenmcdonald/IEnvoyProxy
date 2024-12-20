#!/bin/sh

cd "$(dirname "$0")" || exit 1

cd hysteria || exit 1
git reset --hard
cd ../v2ray-core || exit 1
git reset --hard
git clean -fd
cd ..

patch --directory="hysteria" --strip=1 < hysteria.patch
patch --directory="v2ray-core" --strip=1 < v2ray-core.patch
