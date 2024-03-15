# Use gdate on macOS (GNU coreutils), date on Linux
DATE_CMD="date"
if command -v gdate > /dev/null; then
    DATE_CMD="gdate"
fi

# Mimic preexec functionality using DEBUG trap
preexec_invoke_exec() {
    if [[ "$BASH_COMMAND" != "precmd_invoke_cmd" && "$BASH_COMMAND" != "$PROMPT_COMMAND" ]]; then
        # Save the start time and the command only if it's not the precmd function
        export COMMAND_START_TIME=$($DATE_CMD +%s%3N)
        export LAST_COMMAND="$BASH_COMMAND"
    fi
}
trap 'preexec_invoke_exec' DEBUG

# Mimic precmd functionality using PROMPT_COMMAND
precmd_invoke_cmd() {
    local end_time=$($DATE_CMD +%s%3N)
    local duration=$((end_time - COMMAND_START_TIME))
    {{.CommandScriptPath}} "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
}

PROMPT_COMMAND="precmd_invoke_cmd"
