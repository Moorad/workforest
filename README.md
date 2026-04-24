# Workforest

A simple CLI tool for managing multiple git worktree + tmux sessions

## Todo

- [x] `workforest version` - Show current workforest version and check if dependencies (git, fzf and tmux) are installed
- [] `workforest.yml` - Parse config file and create tmux sessions based on it
  - [] `windows` - array of windows to open (key = window name, value = command to run in window. Alternatively, can be `windows.[].command`)
  - [] `windows.[].include` - Array of ignored to copy when creating a new worktree
  - [] `windows.[].pre-exit` - a hook to run on the window before exiting (or switching) session.
- [] `worktree` - Create tmux session based on workforest file, pick a worktree if multiple were found
- [] `workforest switch` - Fuzzy find tmux sessions
- [] `workforest add` - Create new worktrees in parent directory
- [] `workforest remove` - Delete worktrees in current or sibling directory
- [] `workforest`
