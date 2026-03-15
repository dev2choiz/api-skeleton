#!/bin/bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
SESSION="app"
FIRST_WINDOW="stack"
SELECTED_WINDOW=$FIRST_WINDOW
SELECTED_PANE="0"

tmux kill-session -t $SESSION 2>/dev/null || true

# $1: name of the window to create
create_window() {
  local name=$1

  if ! tmux has-session -t "$SESSION" 2>/dev/null; then
    tmux new-session -d -x "$(tput cols)" -y "$(tput lines)" -s $SESSION -n "$name"
  else
    tmux new-window -n "$name" -t "$SESSION" -d
  fi
}

# $1: pane target (ex: "dev:stack.0")
# $2: command
# $3: validate command = true
send_cmd() {
  local pane=$1
  local cmd=$2
  local execute=${3:-true}

  if [ "$execute" = true ]; then
    tmux send-keys -t "$pane" "$cmd" C-m
  else
    tmux send-keys -t "$pane" "$cmd"
  fi
}

CURR_WINDOW=$FIRST_WINDOW
create_window $CURR_WINDOW
tmux split-window -v -t "$SESSION:$CURR_WINDOW.0" -l 80%
tmux split-window -h -t "$SESSION:$CURR_WINDOW.0"
tmux split-window -h -t "$SESSION:$CURR_WINDOW.2" -l 30%
tmux split-window -v -t "$SESSION:$CURR_WINDOW.3"

send_cmd "$SESSION:$CURR_WINDOW.0" "cd $PROJECT_DIR"
send_cmd "$SESSION:$CURR_WINDOW.0" "make start"
send_cmd "$SESSION:$CURR_WINDOW.1" "cd $PROJECT_DIR"
send_cmd "$SESSION:$CURR_WINDOW.2" "cd $PROJECT_DIR"
send_cmd "$SESSION:$CURR_WINDOW.2" "nvim ."
send_cmd "$SESSION:$CURR_WINDOW.3" "cd $PROJECT_DIR"
send_cmd "$SESSION:$CURR_WINDOW.4" "cd $PROJECT_DIR"

# Attach
tmux select-window -t "$SESSION:$SELECTED_WINDOW"
tmux select-pane -t "$SESSION:$SELECTED_WINDOW.$SELECTED_PANE"
tmux attach -t $SESSION
