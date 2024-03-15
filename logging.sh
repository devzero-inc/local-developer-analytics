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
SOCKET_PATH="/tmp/myapp.socket"

# Determine the system's 'date' command flavor and format the date accordingly
if date --version >/dev/null 2>&1; then
    # GNU 'date'
    FORMATTED_START_DATE=$(date -d "@$4" '+%Y-%m-%d %H:%M:%S')
    FORMATTED_END_DATE=$(date -d "@$5" '+%Y-%m-%d %H:%M:%S')
else
    # BSD 'date'
    FORMATTED_START_DATE=$(date -r "$4" '+%Y-%m-%d %H:%M:%S')
    FORMATTED_END_DATE=$(date -r "$5" '+%Y-%m-%d %H:%M:%S')
fi

# Construct the log message with start time, end time, and duration
LOG_MESSAGE="$2:$3:Start $FORMATTED_START_DATE:End $FORMATTED_END_DATE:Duration $6 seconds:$1"

# Send the log message to the Go application via UNIX socket
echo "$LOG_MESSAGE" | nc -U "$SOCKET_PATH"
