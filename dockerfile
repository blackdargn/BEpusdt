# BEpusdt Docker 镜像
#
# 方式一（默认 target: runtime）：Docker 内完整构建 web + Go（与原 Dockerfile 行为一致）
#   docker compose build
#   docker build -t bepusdt:local .
#
# 方式二（target: runtime-prebuilt）：本地先构建 Linux 二进制后再打镜像
#   scripts/build.sh  或  scripts/build.ps1
#   docker build --target runtime-prebuilt -t bepusdt:local .

# ---------------------------------------------------------------------------
# Stage 1: 构建 Web 管理后台 → static/secure
# ---------------------------------------------------------------------------
FROM node:25.2.1 AS web_builder

RUN npm install -g pnpm@10

WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile --shamefully-hoist

COPY web/ ./
RUN pnpm run build:prod

# ---------------------------------------------------------------------------
# Stage 2: 构建 Go 二进制（嵌入 static/secure 与 static/payment）
# ---------------------------------------------------------------------------
FROM golang:1.26.2-alpine3.23 AS builder

ENV GO111MODULE=on
WORKDIR /go/release

ADD . .

COPY --from=web_builder /web/dist ./static/secure

ARG VERSION=unknown

RUN set -x \
    && MODULE_PATH=$(go list -m) \
    && CGO_ENABLED=0 go build -trimpath \
    -ldflags="-X '${MODULE_PATH}/app.Version=${VERSION}' -s -w -buildid=" \
    -o bepusdt ./main

# ---------------------------------------------------------------------------
# Stage 3a: 运行时（默认，完整构建）
# ---------------------------------------------------------------------------
FROM alpine:3.20 AS runtime

ENV TZ=Asia/Shanghai

RUN apk add --no-cache tzdata ca-certificates \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=builder /go/release/bepusdt /usr/local/bin/bepusdt
COPY static/payment /app/static/payment

WORKDIR /var/lib/bepusdt

EXPOSE 8080

ENTRYPOINT ["bepusdt"]
CMD ["start"]

# ---------------------------------------------------------------------------
# Stage 3b: 运行时（本地预编译二进制，跳过 Docker 内编译）
# ---------------------------------------------------------------------------
FROM alpine:3.20 AS runtime-prebuilt

ENV TZ=Asia/Shanghai

RUN apk add --no-cache tzdata ca-certificates \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY bepusdt /usr/local/bin/bepusdt
COPY static/payment /app/static/payment

WORKDIR /var/lib/bepusdt

EXPOSE 8080

ENTRYPOINT ["bepusdt"]
CMD ["start"]
