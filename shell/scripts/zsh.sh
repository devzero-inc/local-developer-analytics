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
  # Send an end execution message
  {{.CommandScriptPath}} "end" "$LAST_COMMAND" "$PWD" "$USER" "$UUID"
}
