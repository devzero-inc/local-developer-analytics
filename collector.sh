#!/bin/sh
# POSIX-compliant script to log command details to a Unix socket

# Parameters:
# $1 - Command to log
# $2 - Working directory
# $3 - User who executed the command
# $4 - Timestamp of command execution (start time, Unix timestamp)
# $5 - End time of command execution (Unix timestamp)
# $6 - Duration of command execution in seconds

# UNIX socket path
SOCKET_PATH="/tmp/lda.socket"

# Construct the log message with start time, end time, and duration
LOG_MESSAGE="$1|$2|$3|$4|$5|$6"

# Function to check command existence
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Send the log message to the Go application via UNIX socket
if command_exists nc; then
    echo "$LOG_MESSAGE" | nc -U "$SOCKET_PATH"
elif command_exists socat; then
    echo "$LOG_MESSAGE" | socat - UNIX-CONNECT:"$SOCKET_PATH"
else
    echo "Neither nc nor socat is available on this system."
    exit 1
fi
