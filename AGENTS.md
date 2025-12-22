# Agent Coding Guidelines

## Build/Lint/Test Commands
- `make build`: Build binary to `bin/mysql-client-gui`
- `make run`: Run app from source
- `make test`: Run all tests with `go test ./...`
- `go test ./internal/...`: Run tests for specific package
- `go run ./cmd/mysql-client-gui`: Direct run without make
- `CGO_ENABLED=0 go build -trimpath ./cmd/mysql-client-gui`: Static build

## Code Style Guidelines
- Go 1.25+; keep functions small and explicit
- Use stdlib; avoid new runtime dependencies unless necessary
- Naming: `mixedCaps` for locals/functions, `PascalCase` for exported, `UPPER_SNAKE_CASE` for constants
- Import order: stdlib, third-party, internal packages (no grouping required)
- Error handling: return errors with context, no panic()
- Chinese UI messages; keep consistent with existing Chinese prompts/errors
- Use tview for TUI, go-sql-driver/mysql for DB

## Security & Safety
- Never log or persist credentials; keep passwords in memory only
- Never commit database dumps, exported CSVs, or real connection details
- Validate all user inputs; use Config.Validate() for DB configs

## Project Structure
- `cmd/mysql-client-gui/`: App entrypoint
- `internal/`: Internal packages (db, export, tui)
- No test suite yet; add Go testing package tests when needed

## Commit Style
- Conventional Commits (feat:, fix:, refactor:, etc.)
- Git history minimal; single initial commit