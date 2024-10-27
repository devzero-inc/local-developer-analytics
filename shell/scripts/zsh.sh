generate_uuid() {
  echo "$(date +%s)-$$-$RANDOM"
}

preexec() {
  export LAST_COMMAND=$1
  UUID=$(generate_uuid)
  # Send a start execution message
  {{.CommandScriptPath}} "start" "$LAST_COMMAND" "$PWD" "$USER" "$UUID"
}

precmd() {
  local exit_status=$?
  local result="success"
  
  if [[ $exit_status -ne 0 ]]; then
    result="failure"
  fi
  
  # Send an end execution message with result and exit status
  {{.CommandScriptPath}} "end" "$LAST_COMMAND" "$PWD" "$USER" "$UUID" "$result" "$exit_status"
}
