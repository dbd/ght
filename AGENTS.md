# AGENTS.md

Agent guidelines for working in the `ght` (GitHub Terminal) codebase.

## Project Overview

A terminal UI application for browsing and managing GitHub Pull Requests. Built with Go using the Bubble Tea TUI framework from Charm.

**Key dependencies:**
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/glamour` - Markdown rendering
- `github.com/cli/go-gh/v2` - GitHub CLI integration
- `github.com/spf13/viper` - Configuration management

## Commands

```bash
# Build
go build ./...

# Run main application
go run .

# Run sub-commands (development utilities)
go run ./cmd/pr      # PR diff rendering test
go run ./cmd/query   # API query test

# No tests exist yet
go test ./...        # Returns "no test files" for all packages
```

## Architecture

### Bubble Tea (Elm Architecture)

The app follows the Model-Update-View pattern:

1. **Model** - Application state struct
2. **Update(msg)** - Handles messages, returns new model + commands
3. **View()** - Renders model to string
4. **Init()** - Returns initial command(s)

Messages flow: User input → `tea.Msg` → `Update()` → State change → `View()` re-render

### Directory Structure

```
├── main.go                    # Entry point, root Model
├── update.go                  # Root Update() handler
├── view.go                    # Root View() + header/footer
├── commands.go                # Command mode (:merge, :refresh)
├── components/
│   ├── types.go               # Page interface, Context, command types
│   ├── keys.go                # KeyMap with vim-style bindings
│   ├── styles.go              # Lipgloss styles and colors
│   ├── config.go              # Viper config loading
│   ├── render.go              # Box rendering utilities
│   ├── overlay.go             # PlaceOverlay for dialogs
│   ├── mergeDialog.go         # PR merge dialog component
│   ├── pullRequestSearch/     # PR list/search component
│   ├── pullRequestDetail/     # PR detail view + diff
│   └── tab/                   # Tab wrapper component
├── internal/api/
│   ├── pullrequest.go         # GraphQL queries via go-gh
│   ├── types.go               # GitHub API response types
│   └── utils.go               # Helper methods
├── utils/
│   └── gitDiffParse.go        # Git diff text parser
└── cmd/
    ├── pr/                    # Standalone diff render test
    └── query/                 # Standalone API query test
```

### Key Interfaces

**Page interface** (`components/types.go:20`):
```go
type Page interface {
    Init() tea.Cmd
    Update(tea.Msg) (Page, tea.Cmd)
    View() string
    Blur() tea.Msg
    Focus() tea.Msg
    ToggleHelp() tea.Msg
    IsInTextInput() bool
}
```

All major UI components implement this interface.

**Context** (`components/types.go:10`):
Shared state passed to all components:
- Viewport dimensions
- Status text
- KeyMap
- Help model

### Component Hierarchy

```
Model (main.go)
├── []tab.Model
│   ├── pullRequestSearch.Model  (for saved searches)
│   └── pullRequestDetail.Model  (opened PRs)
└── viewport.Model (scrolling)
```

## Key Patterns

### Message-Based Communication

Components communicate via typed messages:
- `pullRequestSearch.OpenPR{PR}` - Open a PR detail tab
- `components.Blur` - Return focus to parent
- `components.CmdMerge{}` - Execute merge command
- `api.PullRequests{}` - API response with PR data

### Keybindings (Vim-style)

Defined in `components/keys.go`:
- `j/k` or arrows - Navigate
- `h/l` or arrows - Switch tabs
- `Enter` - Select/Open
- `/` - Filter
- `?` - Toggle help
- `q` / `Ctrl+W` - Close tab
- `Esc` / `Ctrl+C` - Exit
- `:` - Command mode (leader key)

### Focus Management

- `m.focused` tracks whether root model has focus vs active tab
- Components implement `Blur()`/`Focus()` to manage their focus state
- When at top of a component, `Up` key triggers `Blur()` to return to tab bar

### Styling

Colors defined in `components/styles.go`:
- Green: approvals, additions
- Red: deletions, change requests
- Blue: borders, UI elements
- Yellow: comments
- DarkGrey: background/muted text

Box rendering uses custom Unicode borders (rounded corners: `╭╮╰╯`).

### Configuration

Config file: `~/.config/ght/config.yaml`

```yaml
pr:
  searches:
    - name: "My PRs"
      query: "is:pr author:@me"
    - name: "Review Requested"
      query: "is:pr review-requested:@me"
```

Loaded via Viper in `components/config.go`.

## GitHub API

Uses `go-gh` for GitHub API access (GraphQL):
- Relies on `gh` CLI authentication
- GraphQL queries defined inline in `internal/api/pullrequest.go`
- Also uses `gh pr diff` command for fetching diffs

API types mirror GitHub GraphQL schema with `graphql:"..."` tags.

## Gotchas

1. **Viewport management**: Components must track viewport dimensions from Context and resize appropriately on `tea.WindowSizeMsg`.

2. **Command mode**: Leader key (`:`) only activates when no component is in text input mode (`IsInTextInput()`).

3. **Tab management**: Tabs are stored in a slice. Closing requires rebuilding the slice without the target index.

4. **Overlay rendering**: `PlaceOverlay()` in `overlay.go` handles rendering dialogs/modals on top of existing content.

5. **Diff parsing**: Custom diff parser in `utils/gitDiffParse.go` parses unified diff format. Watch for edge cases with hunk headers.

6. **No tests**: The codebase currently has no test files. Consider adding tests for:
   - Diff parsing (`utils/gitDiffParse.go`)
   - API type parsing
   - Timestamp formatting

7. **GraphQL fragments**: API types use GraphQL union type fragments (`... on PullRequest`). Changes to queries require matching type changes.

## Adding New Features

### New Component
1. Create package under `components/`
2. Implement `Page` interface
3. Add message types for communication
4. Handle keyboard input in `Update()`
5. Use `Context` for viewport/keymap access

### New Command
1. Add command type to `components/types.go`
2. Add to `cmdMap` in `commands.go`
3. Handle in `sendCommandMessage()`
4. Process in appropriate component's `Update()`

### New API Query
1. Define response types in `internal/api/types.go`
2. Add query function in `internal/api/pullrequest.go`
3. Create `tea.Cmd` wrapper returning `tea.Msg`
4. Handle message in component's `Update()`
