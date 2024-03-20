#!/bin/sh

cd "$(dirname "$0")" || exit 1

cd lyrebird || exit 1
git reset --hard
cd ../v2ray-core || exit 1
git reset --hard
git clean -fd
cd ../snowflake || exit 1
git reset --hard
cd ..

patch --directory="lyrebird" --strip=1 < lyrebird.patch
patch --directory="v2ray-core" --strip=1 < v2ray-core.patch
patch --directory="snowflake" --strip=1 < snowflake.patch
