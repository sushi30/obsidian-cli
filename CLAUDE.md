# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Build for local development (creates binary in current directory)
go build -o obsidian-cli

# Build for development with debug name
go build -o obsidian-cli-dev

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build for all platforms (creates bin/darwin/, bin/linux/, bin/windows/)
make build-all
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/actions
go test ./pkg/obsidian

# Run a specific test
go test ./pkg/actions -run TestCreateNote

# Run with verbose output
go test -v ./pkg/actions

# Run with coverage
go test -coverprofile=coverage.out ./...
```

## Architecture

### Core Layered Structure

The codebase follows a three-layer architecture:

1. **cmd/** - CLI commands (Cobra commands)
   - Each command maps to a user-facing CLI operation
   - Commands validate inputs and call functions in the actions layer
   - Command files handle vault name resolution and daily note (@daily) expansion

2. **pkg/actions/** - Business logic layer
   - Orchestrates vault and note operations
   - Handles file I/O, validation, and error handling
   - All actions work with the VaultManager and NoteManager interfaces
   - Functions take parameter structs (e.g., `CreateParams`, `EditParams`)

3. **pkg/obsidian/** - Domain model and interfaces
   - Defines core interfaces: `VaultManager`, `NoteManager`, `UriManager`
   - Implements Obsidian-specific logic (URI construction, path handling, daily notes)
   - Manages the Obsidian config (~/.config/obsidian-cli/config.json)

### Key Components

**Vault Management** (`pkg/obsidian/vault.go`)
- `Vault` struct stores vault name
- `VaultManager` interface provides methods to get vault path, default name, daily note pattern
- Reads from Obsidian's config (~/.config/Obsidian/obsidian.json) and CLI config (~/.config/obsidian-cli/config.json)

**Note Management** (`pkg/obsidian/note.go`)
- `Note` struct implements `NoteManager` interface
- Handles file operations (move, delete, get/set contents)
- Implements link updates when notes are renamed (updates all wikilinks [[]] in vault)
- Supports content search with snippets and backlink finding

**URI Protocol** (`pkg/obsidian/uri.go`)
- Constructs `obsidian://` URIs to open notes in Obsidian app
- `Uri.Construct()` builds URIs with query parameters
- `Uri.Execute()` opens URIs using system default handler

**Daily Notes** (`pkg/obsidian/daily.go`)
- `@daily` is a special reference that resolves to today's daily note
- `ExpandDatePattern()` supports date tokens: YYYY, YY, MM, DD, MMM, MMMM
- Daily note pattern is stored in CLI config and defaults to "YYYY-MM-DD"

**Frontmatter** (`pkg/frontmatter/frontmatter.go`)
- Parses and manipulates YAML frontmatter in markdown files
- Uses `---` delimiters for frontmatter blocks
- Supports setting, deleting, and reading frontmatter keys

### Testing with Mocks

The `mocks/` directory contains test doubles for interfaces:
- `MockVaultOperator` implements `VaultManager`
- `MockNoteManager` implements `NoteManager`
- `MockUriManager` implements `UriManager`

Tests use these mocks to avoid file system operations. Tests that need actual files use `t.TempDir()`.

### Special Patterns

**@daily Note Reference**: The string `@daily` is a special note name that resolves to today's daily note path. Commands handle this via `ResolveNoteName()` in `cmd/daily_resolver.go`.

**Path Handling**: Note paths can be:
- Simple names: `"note"` → searches vault for `note.md`
- Relative paths from vault root: `"folder/note"` → `<vault>/folder/note.md`
- Full matches have priority over basename matches

**Editor Integration**: Commands support `--editor` flag to open notes in $EDITOR instead of Obsidian app. GUI editors (code, subl, atom, etc.) automatically get `--wait` flag added.

## Dependencies

- **cobra** - CLI framework
- **go-fuzzyfinder** - Interactive fuzzy search
- **open-golang** - Cross-platform file/URI opening
- **adrg/frontmatter** - YAML frontmatter parsing
- **stretchr/testify** - Testing assertions and mocks

## Version and Release

Version is stored in `cmd/root.go` as `Version: "vX.Y.Z"`. The Makefile includes targets for automated releases:
- `make release VERSION=vX.Y.Z` - Full release process (update version, build, tag, push)
- `make release-patch/minor/major` - Auto-increment version and release
