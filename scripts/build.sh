#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VERSION="${VERSION:-dev-local}"

echo "==> 安装并构建 Web 前端..."
npm install -g pnpm@10
pushd web >/dev/null
pnpm install --frozen-lockfile --shamefully-hoist
pnpm run build:prod
popd >/dev/null

echo "==> 同步前端产物到 static/secure..."
mkdir -p static/secure
rm -rf static/secure/*
cp -r web/dist/* static/secure/

echo "==> 构建 Linux 二进制（嵌入 static）..."
MODULE_PATH="$(go list -m)"
CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH:-amd64}" go build -trimpath \
  -ldflags="-X '${MODULE_PATH}/app.Version=${VERSION}' -s -w -buildid=" \
  -o bepusdt ./main

echo ""
echo "构建完成: ${ROOT}/bepusdt"
echo "打包镜像: docker build --target runtime-prebuilt -t bepusdt:local ."
