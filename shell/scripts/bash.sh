generate_uuid() {
    echo "$(date +%s)-$$-$RANDOM" | shasum | cut -d " " -f1
}

preexec_invoke_exec() {
# Mimic preexec functionality using DEBUG trap
    # Avoid running preexec_invoke_exec for PROMPT_COMMAND
    if [[ "$BASH_COMMAND" != "$PROMPT_COMMAND" ]]; then
        export UUID=$(generate_uuid)
        export LAST_COMMAND="$BASH_COMMAND"
        # Send a start execution message
        {{.CommandScriptPath}} "start" "$LAST_COMMAND" "$PWD" "$USER" "$UUID"
    fi
}
trap 'preexec_invoke_exec' DEBUG

# Mimic precmd functionality using PROMPT_COMMAND
precmd_invoke_cmd() {
    # Send an end execution message
    {{.CommandScriptPath}} "end" "$LAST_COMMAND" "$PWD" "$USER" "$UUID"
}

# Update PROMPT_COMMAND to invoke precmd_invoke_cmd
# Append precmd_invoke_cmd to existing PROMPT_COMMAND to preserve other functionalities
PROMPT_COMMAND="${PROMPT_COMMAND:+$PROMPT_COMMAND; }precmd_invoke_cmd"