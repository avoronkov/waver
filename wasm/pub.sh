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

env GOOS=js GOARCH=wasm go build -o ./static/assets/demo.wasm ./cmd/wasm
gzip -f -9 ./static/assets/demo.wasm
cp -r ./static/assets ../public

echo 'OK'
