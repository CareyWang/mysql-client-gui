# Repository Guidelines

## Project Structure & Module Organization
- `cmd/mysql-client-gui/`: Go app entrypoint (tview TUI).
- `internal/`: Internal packages such as `db`, `export`, and `tui`.
- `scripts/`: Helper scripts for local build/run of packaged app.
- `dist/`, `build/`: Generated build outputs (do not edit by hand).
- `README.md`: User-facing run and packaging instructions (Linux-focused).

## Build, Test, and Development Commands
- `go run ./cmd/mysql-client-gui`: Run the app from source.
- `go test ./...`: Build and run tests (currently a basic smoke check).
- `./scripts/build-linux-onedir.sh`: Build a Linux binary into `dist/mysql-client-gui/` (defaults to `linux/amd64`).
- `./scripts/run-linux-onedir.sh`: Run the packaged binary from `dist/`.
- Docker (non-Linux hosts):
  - `docker build -f Dockerfile.linux-build -t mysql-client-gui:linux-build .`
  - `docker run --rm -v "$PWD/dist:/out" mysql-client-gui:linux-build`

## Coding Style & Naming Conventions
- Go 1.22+; keep functions small and explicit.
- Prefer stdlib; avoid new runtime dependencies unless necessary.
- Naming: `mixedCaps` for locals/functions, `PascalCase` for exported identifiers, constants in `UPPER_SNAKE_CASE`.
- UI/user messages should stay consistent with existing language (many prompts/errors are Chinese).

## Testing Guidelines
- No dedicated test suite yet; use Go’s `testing` package under `internal/...` or `cmd/...` when adding tests.
- Minimal smoke checks before PRs: `go test ./...` and run the UI via `go run ./cmd/mysql-client-gui` to verify connect/query/export flows.

## Commit & Pull Request Guidelines
- Git history is minimal (single initial commit). Recommended format: Conventional Commits (e.g., `feat: add CSV export option`).
- PRs should include a clear description, verification steps, and screenshots for UI changes. Note OS/desktop environment if relevant (Linux/X11/Wayland).

## Security & Configuration Tips
- Do not log or persist credentials; keep passwords in memory only.
- Never commit database dumps, exported CSVs, or real connection details.
