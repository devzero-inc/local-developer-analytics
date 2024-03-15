# LDA

## Adding scripts

Add `logging.sh` to your system

Run `chmod + x logging.sh` to make script executable

Add script to your shell, this is dependant on what kind of shell are you using
for zsh use:

```
preexec() {
    export COMMAND_START_TIME=$(date +%s)
    # Optionally, capture the command itself if needed
    export LAST_COMMAND=$1
}

precmd() {
    local end_time=$(date +%s)
    local duration=$((end_time - COMMAND_START_TIME))
    # Now call the logging script with additional end time and duration
    $DOT/zsh/logging.sh "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
}

```
