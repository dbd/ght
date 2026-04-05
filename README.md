# ght - GitHub TUI

A fast, keyboard-driven terminal UI for browsing and managing GitHub Pull Requests, Issues, and Milestones.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the [GitHub CLI](https://cli.github.com/).

## Features

- 🔍 **Powerful Search** - Search PRs and issues with GitHub's query syntax
- 📑 **Multiple Tabs** - Browse multiple PRs, issues, and searches simultaneously
- ⌨️ **Vim-style Keys** - Navigate efficiently with familiar keybindings
- 📝 **Full PR Management** - Comment, approve, request changes, merge
- 🐛 **Issue Management** - View, comment, close/reopen, and assign issues
- 🏁 **Milestones** - Browse milestones with progress bars and drill into their issues
- 👥 **Team Collaboration** - Add reviewers and assignees
- 💾 **Saved Searches** - Save frequently used PR and issue searches to config
- 🎨 **Rich Rendering** - Markdown support, syntax highlighting, colored diffs
- 🔄 **Live Updates** - Refresh PRs and searches on demand

## Installation

### Prerequisites

- [GitHub CLI (`gh`)](https://cli.github.com/) - Must be installed and authenticated
- Go 1.21+ (for installation)

### Using go install (Recommended)

```bash
go install github.com/dbd/ght@latest
```

This installs the `ght` binary to `$GOPATH/bin` (usually `~/go/bin`).

**Standard Go installation:**
Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH="$HOME/go/bin:$PATH"
```

**If using asdf:**
Add asdf's Go packages bin to your PATH:
```bash
export PATH="$HOME/.asdf/installs/golang/$(asdf current golang | awk '{print $2}')/packages/bin:$PATH"
```

### From Source

```bash
git clone https://github.com/dbd/ght.git
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
   issue:
     searches:
       - name: "My Issues"
         query: "is:issue author:@me is:open"
       - name: "Assigned to Me"
         query: "is:issue assignee:@me is:open"
     milestones:
       - name: "My Project"
         repo: "owner/repo"
   ```

3. **Launch ght**:
   ```bash
   ./ght
   ```

## Usage

### Modes

ght has two modes: **PR mode** and **Issue mode**. Switch between them with commands:

| Command | Description |
|---------|-------------|
| `:prs` | Switch to PR mode |
| `:issues` | Switch to Issue mode |

### Navigation

| Key | Action |
|-----|--------|
| `j`/`k` or `↑`/`↓` | Move up/down |
| `h`/`l` or `←`/`→` | Switch tabs left/right |
| `Enter` | Open selected item |
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
| `A` | Add assignee |
| `m` | Open merge dialog |

### Issue Actions

| Key | Action |
|-----|--------|
| `c` | Add comment |
| `A` | Add assignee |
| `x` | Close / reopen issue |
| `o` | Open issue in browser |
| `M` | Open milestone (if issue has one) |

### Command Mode

Press `:` to enter command mode. Available commands:

| Command | Description |
|---------|-------------|
| `:prs` | Switch to PR mode |
| `:issues` | Switch to Issue mode |
| `:newtab` | Create a new PR search tab |
| `:new-issue-tab` | Create a new issue search tab |
| `:milestones <owner/repo>` | Open milestone list for a repo |
| `:save-tab <name>` | Save current search to config |
| `:refresh` | Refresh current tab |
| `:merge` | Open merge dialog |
| `:add-assignee <username>` | Add assignee to PR |
| `:add-reviewer <username>` | Add reviewer to PR |
| `:comment <message>` | Add a comment |
| `:approve [message]` | Approve PR with optional comment |
| `:request-changes <message>` | Request changes on PR |
| `:help` | Show help dialog |
| `:quit` | Quit ght |

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

issue:
  searches:
    - name: "My Issues"
      query: "is:issue author:@me is:open"

    - name: "Assigned to Me"
      query: "is:issue assignee:@me is:open"

    - name: "Bugs"
      query: "is:issue label:bug is:open"

  milestones:
    - name: "My Project"
      repo: "owner/repo"
    - name: "Other Repo"
      repo: "owner/other-repo"
```

### GitHub Search Query Syntax

Use GitHub's advanced search syntax for powerful queries:

- `is:pr` / `is:issue` - Filter by type
- `is:open` / `is:closed` / `is:merged` - State filter
- `author:username` / `author:@me` - Filter by author
- `review-requested:@me` - PRs requesting your review
- `assignee:username` - Assigned to user
- `org:orgname` / `repo:owner/repo` - Scope to org or repo
- `label:bug` - Filter by label
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

### Issue Triage

1. Switch to issue mode with `:issues`
2. Browse issues with `j`/`k`, press `Enter` to open
3. Press `x` to close/reopen, `A` to add an assignee, `c` to comment
4. Press `M` to jump to the issue's milestone

### Milestone Tracking

1. Open a milestone list: `:milestones owner/repo`
2. Browse milestones — progress bars show open/closed issue counts
3. Press `Enter` to drill into a milestone's issues
4. Press `Enter` on an issue to open it

### Quick Comment

From any PR or issue detail view:
```
:comment Great work! LGTM
```

Or use the dialog:
1. Press `C` (PR) or `c` (issue)
2. Type your comment
3. Press `Enter` / `Ctrl+S` to submit

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
│   ├── issueSearch/         # Issue list/search component
│   ├── issueDetail/         # Issue detail view
│   ├── milestoneList/       # Milestone list component
│   ├── milestoneDetail/     # Milestone detail + issue table
│   └── tab/                 # Tab wrapper component
├── internal/api/
│   ├── pullrequest.go       # PR GitHub API calls
│   ├── issue.go             # Issue GitHub API calls
│   ├── milestone.go         # Milestone GitHub API calls
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

**PR search:** Use `:newtab` to create a new search tab, enter your query with `/`, then save it with `:save-tab "My Custom Search"`.

**Issue search:** Use `:new-issue-tab`, enter your query with `/`, then save with `:save-tab "My Issue Search"`.

### Quick Access

In any list, press `/` and type to filter by title, number, or other fields.

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
