# Workforest

A simple CLI tool for managing multiple git worktree + tmux sessions

## Todo

- [x] `workforest version` - Show current workforest version and check if dependencies (git and tmux) are installed
- [x] non-worktree `worktree` - Create tmux session based on workforest file
- [x] worktree `workforest` - Pick a worktree (if multiple were found)
- [] `workforest.yml` - Parse config file and create tmux sessions based on it
  - [x] `windows` - array of windows to open
  - [] `windows.[].pre-exit` - a hook to run on the window before exiting (or switching) session.
- [] `workforest switch` - Fuzzy find tmux sessions
- [] `workforest add` - Create new worktrees in parent directory
  - [] `windows.[].include` - Array of ignored to copy when creating a new worktree
- [] `workforest remove` - Delete worktrees in current or sibling directory
- [] Better error handling and error messages
