#!/bin/sh

set -ex

base_dir() {
	dirname $(readlink -f $0)
}

die() {
	echo "$@" >&2
	exit 1
}

BASE_DIR="$(base_dir)"
cd "$BASE_DIR" || dir "Failed to chdir '$BASE_DIR'"

mkdir -p ./waver/web
cp -rT ./web/ ./waver/web
env GOOS=js GOARCH=wasm go build -o ./waver/web/app.wasm ./cmd/pwa
go run ./cmd/pages
rm -rf ../public
cp -rT ./waver ../public

echo 'OK'
