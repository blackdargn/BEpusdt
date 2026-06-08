$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

$Version = if ($env:VERSION) { $env:VERSION } else { "dev-local" }
$GoArch = if ($env:GOARCH) { $env:GOARCH } else { "amd64" }

Write-Host "==> 安装并构建 Web 前端..."
npm install -g pnpm@10
Push-Location web
pnpm install --frozen-lockfile --shamefully-hoist
pnpm run build:prod
Pop-Location

Write-Host "==> 同步前端产物到 static/secure..."
New-Item -ItemType Directory -Force -Path static/secure | Out-Null
Remove-Item -Recurse -Force static/secure/* -ErrorAction SilentlyContinue
Copy-Item -Recurse -Force web/dist/* static/secure/

Write-Host "==> 构建 Linux 二进制（嵌入 static）..."
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = $GoArch
$ModulePath = go list -m
go build -trimpath `
  -ldflags="-X '${ModulePath}/app.Version=${Version}' -s -w -buildid=" `
  -o bepusdt ./main

Write-Host ""
Write-Host "构建完成: $Root\bepusdt"
Write-Host "打包镜像: docker build --target runtime-prebuilt -t bepusdt:local ."
