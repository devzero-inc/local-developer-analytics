# Use gdate on macOS (GNU coreutils), date on Linux
if type -q gdate
    set DATE_CMD gdate
else
    set DATE_CMD date
end

function preexec --on-event fish_preexec
    # Capture the command start time in milliseconds
    set -gx COMMAND_START_TIME ($DATE_CMD +%s%3N)
    set -gx LAST_COMMAND $argv
end

function precmd --on-event fish_prompt
    set end_time ($DATE_CMD +%s%3N)
    set duration (math $end_time - $COMMAND_START_TIME)
    # Call the logging script with command execution details
    {{.CommandScriptPath}} "$LAST_COMMAND" "$PWD" "$USER" "$COMMAND_START_TIME" "$end_time" "$duration"
end
