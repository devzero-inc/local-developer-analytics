#!/bin/bash

# Function to send data to the Go daemon
send_to_daemon() {
    local end_time=$(date +%s%N)
    local duration=$(( (end_time - COMMAND_START_TIME) / 1000000 ))
    local command_info="$PWD, $LAST_COMMAND, $duration ms"
    echo "$command_info" | nc localhost 8080 &
}

# Initialize the start time and last command
COMMAND_START_TIME=$(date +%s%N)
LAST_COMMAND=""

# Trap DEBUG to capture every command
trap 'LAST_COMMAND=$BASH_COMMAND; COMMAND_START_TIME=$(date +%s%N)' DEBUG

# Function to be executed after a command is run
postexec() {
    send_to_daemon
}

# Set PROMPT_COMMAND to execute after each command
PROMPT_COMMAND=postexec
