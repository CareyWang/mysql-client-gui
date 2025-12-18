# Repository Guidelines

## Project Structure

- `mysql_client_gui/`: Python package (Tkinter GUI + wrappers around `mysql`/`mysqldump`).
- `scripts/`: local helper scripts (build + run packaged app).
- `packaging/pyinstaller/`: PyInstaller spec for Linux “onedir” builds.
- `dist/`, `build/`: build outputs (generated; do not edit by hand).
- `README.md`: user-facing run + packaging instructions (Linux-focused).

## Build, Test, and Development Commands

- `python3 -m mysql_client_gui`: run the GUI from source.
- `python3 -m pip install -e .`: install editable CLI entrypoint `mysql-client-gui`.
- `./scripts/build-linux-onedir.sh`: build a Linux onedir bundle via PyInstaller (expects Python ≥ 3.10 and Tkinter).
- `./scripts/run-linux-onedir.sh`: run the packaged binary from `dist/`.
- Docker build (non-Linux hosts): `docker build -f Dockerfile.linux-build -t mysql-client-gui:linux-build .` and `docker run --rm -v "$PWD/dist:/out" mysql-client-gui:linux-build`.

## Coding Style & Naming Conventions

- Python 3.10+ with type hints; use 4-space indentation and keep functions small and explicit.
- Prefer stdlib (Tkinter, `subprocess`, `dataclasses`) and avoid adding new runtime dependencies unless necessary.
- Naming: `snake_case` for functions/vars, `PascalCase` for classes, constants in `UPPER_SNAKE_CASE`.
- Keep UI/user messages consistent with existing language (currently Chinese in many prompts/errors).

## Testing Guidelines

- No automated test suite is currently checked in. If adding tests, prefer `pytest` and place them under `tests/` with names like `test_*.py`.
- Minimal smoke checks before PRs: run `python3 -m mysql_client_gui`, verify connect/test/query/export flows, and run `python -m compileall mysql_client_gui`.

## Commit & Pull Request Guidelines

- Git history is minimal (currently a single “Initial commit…”); no established convention yet.
- Recommended commit format: Conventional Commits (e.g., `feat: add CSV export option`, `fix: handle empty result set`).
- PRs should include: a clear description, reproduction/verification steps, and screenshots for UI changes. Note OS/desktop environment if relevant (Linux/X11/Wayland).

## Security & Configuration Tips

- Do not log or persist credentials. This project uses a temporary `--defaults-extra-file` for `mysql`/`mysqldump`; keep that pattern when modifying connection logic.
- Never commit database dumps, exported CSVs, or real connection details.
