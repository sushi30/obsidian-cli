# obsidian-cli

Go CLI for interacting with Obsidian vaults from the terminal.

## Build & Test

```bash
go build -o obsidian-cli .
go test ./...
```

Unit tests use `t.TempDir()` for isolation — no external state needed.

## Manual / Integration Testing

Use `./tmp/test-vault/` as a scratch vault for manual and integration testing.
This directory is gitignored. Create it as needed:

```bash
mkdir -p tmp/test-vault
```

## Project Structure

- `cmd/` — Cobra command definitions
- `pkg/actions/` — Business logic for each command
- `pkg/obsidian/` — Vault, note, and URI abstractions
- `mocks/` — Test mocks for interfaces
