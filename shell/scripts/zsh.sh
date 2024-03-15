preexec() {
    # Capture the command start time in milliseconds
    export COMMAND_START_TIME=$(gdate +%s%3N)
    export LAST_COMMAND=$1
}

precmd() {
    local end_time=$(gdate +%s%3N)
    local duration=$((end_time - COMMAND_START_TIME))
    # Call the logging script with command execution details
    {{.CommandScriptPath}} "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
}
