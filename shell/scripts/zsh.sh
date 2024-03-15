# Use gdate on macOS (GNU coreutils), date on Linux
DATE_CMD="date"
if command -v gdate > /dev/null; then
    DATE_CMD="gdate"
fi

preexec() {
    # Capture the command start time in milliseconds
    export COMMAND_START_TIME=$($DATE_CMD +%s%3N)
    export LAST_COMMAND=$1
}

precmd() {
    local end_time=$($DATE_CMD +%s%3N)
    local duration=$((end_time - COMMAND_START_TIME))
    # Call the logging script with command execution details
    {{.CommandScriptPath}} "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
}
