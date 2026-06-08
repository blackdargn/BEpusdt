#!/bin/bash

APP_NAME="pay_server.bin"
APP_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$APP_DIR/app.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "[$APP_NAME] PID file not found, may not be running"
    exit 0
fi

PID=$(cat "$PID_FILE")

if ! kill -0 "$PID" 2>/dev/null; then
    echo "[$APP_NAME] process (PID: $PID) not found, cleaning up PID file..."
    rm -f "$PID_FILE"
    exit 0
fi

echo "[$APP_NAME] stopping (PID: $PID)..."
kill -SIGTERM "$PID"

for i in $(seq 1 10); do
    if ! kill -0 "$PID" 2>/dev/null; then
        rm -f "$PID_FILE"
        echo "[$APP_NAME] stopped successfully"
        exit 0
    fi
    sleep 1
done

echo "[$APP_NAME] not responding, force killing..."
kill -SIGKILL "$PID"
sleep 1
rm -f "$PID_FILE"
echo "[$APP_NAME] force killed"
