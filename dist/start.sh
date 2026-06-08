#!/bin/bash

APP_NAME="pay_server.bin"
APP_DIR="$(cd "$(dirname "$0")" && pwd)"
APP_BIN="$APP_DIR/$APP_NAME"
PID_FILE="$APP_DIR/app.pid"
LOG_FILE="$APP_DIR/app.log"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "[$APP_NAME] already running (PID: $PID)"
        exit 0
    else
        echo "[$APP_NAME] stale PID file found, cleaning up..."
        rm -f "$PID_FILE"
    fi
fi

if [ ! -f "$APP_BIN" ]; then
    echo "[ERROR] $APP_BIN not found"
    exit 1
fi

chmod +x "$APP_BIN"

cd "$APP_DIR"
nohup "$APP_BIN" >> /dev/null 2>&1 &

PID=$!
echo $PID > "$PID_FILE"

sleep 1
if kill -0 "$PID" 2>/dev/null; then
    echo "[$APP_NAME] started successfully (PID: $PID)"
    echo "Log: $LOG_FILE"
else
    echo "[ERROR] $APP_NAME failed to start, check: $LOG_FILE"
    rm -f "$PID_FILE"
    exit 1
fi
