# ght - GitHub TUI

A fast, keyboard-driven terminal UI for browsing and managing GitHub Pull Requests.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the [GitHub CLI](https://cli.github.com/).

## Features

- 🔍 **Powerful Search** - Search PRs with GitHub's query syntax
- 📑 **Multiple Tabs** - Browse multiple PRs and searches simultaneously
- ⌨️ **Vim-style Keys** - Navigate efficiently with familiar keybindings
- 📝 **Full PR Management** - Comment, approve, request changes, merge
- 👥 **Team Collaboration** - Add reviewers and assignees
- 💾 **Saved Searches** - Save frequently used searches to config
- 🎨 **Rich Rendering** - Markdown support, syntax highlighting, colored diffs
- 🔄 **Live Updates** - Refresh PRs and searches on demand

## Installation

### Prerequisites

- [GitHub CLI (`gh`)](https://cli.github.com/) - Must be installed and authenticated
- Go 1.21+ (for building from source)

### From Source

```bash
git clone https://github.com/yourusername/ght.git
cd ght
go build
./ght
```

## Quick Start

1. **Authenticate with GitHub CLI** (if not already done):
   ```bash
   gh auth login
   ```

2. **Configure your searches** in `~/.config/ght/config.yaml`:
   ```yaml
   pr:
     searches:
       - name: "My PRs"
         query: "is:pr author:@me"
       - name: "Review Queue"
         query: "is:pr review-requested:@me"
       - name: "Team PRs"
         query: "is:pr org:yourorg is:open"
   ```

3. **Launch ght**:
   ```bash
   ./ght
   ```

## Usage

### Navigation

| Key | Action |
|-----|--------|
| `j`/`k` or `↑`/`↓` | Move up/down |
| `h`/`l` or `←`/`→` | Switch tabs left/right |
| `Enter` | Open selected PR |
| `q` or `Ctrl+W` | Close current tab |
| `Esc` or `Ctrl+C` | Exit application / Cancel dialog |
| `?` | Toggle help |
| `Ctrl+Z` | Suspend to shell |

### Search & Filter

| Key | Action |
|-----|--------|
| `/` | Enter search/filter mode |
| `Enter` | Execute search (in search mode) |
| `Esc` | Cancel search |

### PR Actions

| Key | Action |
|-----|--------|
| `c` | Show/hide comments |
| `C` | Add comment |
| `a` | Approve PR |
| `x` | Request changes |
| `r` | Add reviewer |
| `A` | Add assignee (Shift+A) |
| `m` | Open merge dialog |

### Command Mode

Press `:` to enter command mode. Available commands:

| Command | Description |
|---------|-------------|
| `:newtab` | Create a new search tab |
| `:save-tab <name>` | Save current search to config |
| `:refresh` | Refresh current tab |
| `:merge` | Open merge dialog |
| `:add-assignee <username>` | Add assignee to PR |
| `:add-reviewer <username>` | Add reviewer to PR |
| `:comment <message>` | Add a comment to PR |
| `:approve [message]` | Approve PR with optional comment |
| `:request-changes <message>` | Request changes on PR |
| `:help` | Show help dialog |

## Configuration

Configuration file location: `~/.config/ght/config.yaml`

### Example Configuration

```yaml
pr:
  searches:
    - name: "My Open PRs"
      query: "is:pr author:@me is:open"
    
    - name: "Needs My Review"
      query: "is:pr review-requested:@me is:open"
    
    - name: "Recently Updated"
      query: "is:pr is:open sort:updated-desc"
    
    - name: "Assigned to Me"
      query: "is:pr assignee:@me is:open"
```

### GitHub Search Query Syntax

Use GitHub's advanced search syntax for powerful queries:

- `is:pr` - Pull requests only
- `is:open` / `is:closed` / `is:merged` - PR state
- `author:username` - PRs by author
- `author:@me` - Your PRs
- `review-requested:@me` - PRs requesting your review
- `assignee:username` - PRs assigned to user
- `org:orgname` - PRs in organization
- `repo:owner/repo` - PRs in specific repo
- `label:bug` - PRs with label
- `sort:updated-desc` - Sort by update time

[Full search syntax documentation](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests)

## Workflows

### Review Workflow

1. Start `ght` and select "Review Queue" tab
2. Browse PRs with `j`/`k`
3. Press `Enter` to view PR details
4. Read diff and comments
5. Press `a` to approve, or `x` to request changes
6. Enter your review message and press `Ctrl+S`
7. Press `q` to close tab and return to queue

### Quick Comment

From any PR detail view:
```
:comment Great work! LGTM 🚀
```

Or use the dialog:
1. Press `C`
2. Type your comment
3. Press `Ctrl+S` to submit

### Merge PR

1. Open PR detail
2. Press `m` or type `:merge`
3. Select merge method (merge commit, squash, rebase)
4. Confirm with `Ctrl+S`

### Add Reviewer

Quick command:
```
:add-reviewer username
```

Or use dialog:
1. Press `r`
2. Enter username
3. Press `Enter`

## Development

### Project Structure

```
├── main.go                  # Entry point
├── update.go                # Root update handler
├── view.go                  # Root view + header/footer
├── commands.go              # Command mode handlers
├── components/
│   ├── types.go             # Page interface, command types
│   ├── keys.go              # Key bindings
│   ├── styles.go            # Lipgloss styles
│   ├── config.go            # Configuration management
│   ├── mergeDialog.go       # Merge dialog component
│   ├── reviewDialog.go      # Review dialog component
│   ├── inputDialog.go       # Input dialog component
│   ├── helpDialog.go        # Help dialog component
│   ├── pullRequestSearch/   # PR list/search component
│   ├── pullRequestDetail/   # PR detail view + diff
│   └── tab/                 # Tab wrapper component
├── internal/api/
│   ├── pullrequest.go       # GitHub API calls
│   └── types.go             # API response types
└── utils/
    └── gitDiffParse.go      # Git diff parser
```

### Testing Commands

```bash
# Run main application
go run .

# Run dev utilities
go run ./cmd/pr      # PR diff rendering test
go run ./cmd/query   # API query test
```

### Adding Features

See [AGENTS.md](AGENTS.md) for detailed development documentation.

## Tips & Tricks

### Create Custom Searches

Use `:newtab` to create a new search tab, enter your query with `/`, then save it with `:save-tab "My Custom Search"`.

### Quick PR Access

In your PR list, press `/` and type the PR number to quickly filter.

### Terminal Multiplexer Integration

Use with `tmux` or `screen`:
- `Ctrl+Z` suspends ght to shell
- `fg` returns to ght
- Or run in dedicated tmux pane

### Markdown in Comments

Comments support full markdown syntax:
```
:comment ## Summary

This PR looks great!

- [x] Code reviewed
- [x] Tests passing
- [ ] Documentation updated
```

## Troubleshooting

### "gh not found" error

Install and authenticate GitHub CLI:
```bash
# Install (macOS)
brew install gh

# Install (Linux)
# See https://github.com/cli/cli/blob/trunk/docs/install_linux.md

# Authenticate
gh auth login
```

### Configuration not loading

Ensure config file exists:
```bash
mkdir -p ~/.config/ght
touch ~/.config/ght/config.yaml
```

### API rate limiting

GitHub CLI uses your authenticated token. Check rate limits:
```bash
gh api rate_limit
```

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- Uses [Glamour](https://github.com/charmbracelet/glamour) for markdown rendering
- Powered by [GitHub CLI](https://cli.github.com/)
