#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

PYTHON_BIN="${PYTHON_BIN:-python3.14}"
DIST_DIR="${DIST_DIR:-dist}"

if ! command -v "$PYTHON_BIN" >/dev/null 2>&1; then
  echo "ERROR: 找不到 Python：$PYTHON_BIN" >&2
  exit 1
fi

echo "[1/4] 检查 Python 版本"
"$PYTHON_BIN" - <<'PY'
import sys
major, minor = sys.version_info[:2]
if (major, minor) < (3, 10):
    raise SystemExit(f"Python >= 3.10 required, got {sys.version}")
print(sys.version)
PY

echo "提示：请在 CentOS 7 x86_64 上执行该脚本以获得最大兼容性（glibc 2.17）。"

echo "[2/4] 创建打包用 venv（.venv-build）"
rm -rf .venv-build
"$PYTHON_BIN" -m venv .venv-build
source .venv-build/bin/activate
python -m pip install -U pip wheel setuptools

echo "[2.5/4] 检查 tkinter 可用性"
python - <<'PY'
try:
    import tkinter  # noqa: F401
except Exception as e:
    raise SystemExit(
        "ERROR: 当前 Python 缺少 tkinter（通常是缺少 _tkinter）。\n"
        "该项目 GUI 依赖 tkinter，打包前必须先解决。\n"
        "\n"
        "建议：\n"
        "- CentOS 7: sudo yum install -y tk tcl tk-devel tcl-devel\n"
        "- Ubuntu/Debian: sudo apt-get install -y python3-tk\n"
        "- macOS: 请使用带 Tcl/Tk 的 Python（例如 python.org 发行版或正确配置的 Homebrew Python）\n"
    ) from e
PY

echo "[3/4] 安装 PyInstaller"
python -m pip install -U pyinstaller

echo "[4/4] 打包（onedir）"
rm -rf "$DIST_DIR" build
pyinstaller --noconfirm --clean packaging/pyinstaller/mysql-client-gui.spec --distpath "$DIST_DIR"

echo "OK: 输出目录：$DIST_DIR/mysql-client-gui/"
