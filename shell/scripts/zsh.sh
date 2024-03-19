# Use gdate on macOS (GNU coreutils), date on Linux
DATE_CMD="date"
if command -v gdate > /dev/null; then
    DATE_CMD="gdate"
fi

preexec() {
    export COMMAND_START_TIME=$($DATE_CMD +%s%3N)
    export LAST_COMMAND=$1
    # Send a start execution message
    {{.CommandScriptPath}} "start" "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME"
}

precmd() {
    local end_time=$($DATE_CMD +%s%3N)
    local duration=$((end_time - COMMAND_START_TIME))
    # Send an end execution message
    {{.CommandScriptPath}} "end" "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
}
