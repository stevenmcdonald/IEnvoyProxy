#!/bin/sh

cd "$(dirname "$0")" || exit 1

cd hysteria || exit 1
git reset --hard
cd ..

patch --directory="hysteria" --strip=1 < hysteria.patch
