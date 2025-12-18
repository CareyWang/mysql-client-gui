# -*- mode: python ; coding: utf-8 -*-

from __future__ import annotations

from pathlib import Path

from PyInstaller.utils.hooks import collect_submodules

block_cipher = None

project_root = Path(__file__).resolve().parents[2]
entry_script = project_root / "mysql_client_gui" / "__main__.py"

try:
    import tkinter  # noqa: F401
except Exception as e:  # pragma: no cover
    raise SystemExit(
        "ERROR: tkinter 不可用（缺少 _tkinter）。\n"
        "请安装带 tkinter 的 Python，或在系统中安装 tcl/tk 相关依赖后重新安装/编译 Python。\n"
        "CentOS 7 参考：yum install -y tk tcl tk-devel tcl-devel"
    ) from e

hiddenimports: list[str] = []
hiddenimports += collect_submodules("tkinter")

a = Analysis(
    [str(entry_script)],
    pathex=[str(project_root)],
    binaries=[],
    datas=[],
    hiddenimports=hiddenimports,
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name="mysql-client-gui",
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=False,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=True,
    upx_exclude=[],
    name="mysql-client-gui",
)
