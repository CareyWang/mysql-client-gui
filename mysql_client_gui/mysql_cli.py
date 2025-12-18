from __future__ import annotations

import os
import shlex
import stat
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Sequence


@dataclass(frozen=True)
class MySQLClientConfig:
    host: str
    port: int = 3306
    user: str = ""
    password: str = ""
    database: str | None = None
    connect_timeout_sec: int = 5
    charset: str = "utf8mb4"

    def _defaults_file_text(self) -> str:
        lines = ["[client]"]
        lines.append(f"host={self.host}")
        lines.append(f"port={self.port}")
        if self.user:
            lines.append(f"user={self.user}")
        if self.password:
            lines.append(f"password={self.password}")
        if self.database:
            lines.append(f"database={self.database}")
        lines.append(f"default-character-set={self.charset}")
        return "\n".join(lines) + "\n"


@dataclass(frozen=True)
class QueryResult:
    columns: list[str]
    rows: list[list[str]]
    raw_stdout: str


class MySQLCLIError(RuntimeError):
    def __init__(self, message: str, *, command: Sequence[str] | None = None, stderr: str | None = None):
        super().__init__(message)
        self.command = list(command) if command else None
        self.stderr = stderr


def _ensure_executable_exists(exe: str) -> None:
    from shutil import which

    if which(exe) is None:
        raise MySQLCLIError(f"找不到可执行文件：{exe}（请先安装 MySQL CLI 客户端并确保在 PATH 中）")


def _write_temp_defaults_file(text: str) -> Path:
    tmp = tempfile.NamedTemporaryFile(prefix="mysql-client-gui-", suffix=".cnf", delete=False)
    try:
        tmp.write(text.encode("utf-8"))
        tmp.flush()
    finally:
        tmp.close()

    path = Path(tmp.name)
    try:
        os.chmod(path, stat.S_IRUSR | stat.S_IWUSR)  # 0o600
    except OSError:
        pass
    return path


def _run_command(args: list[str], *, timeout_sec: int | None = None) -> subprocess.CompletedProcess[str]:
    try:
        return subprocess.run(
            args,
            text=True,
            capture_output=True,
            timeout=timeout_sec,
            check=False,
        )
    except subprocess.TimeoutExpired as e:
        raise MySQLCLIError(f"命令执行超时（>{timeout_sec}s）：{shlex.join(args)}") from e
    except OSError as e:
        raise MySQLCLIError(f"无法执行命令：{shlex.join(args)}") from e


def run_query(config: MySQLClientConfig, sql: str, *, mysql_exe: str = "mysql", timeout_sec: int = 60) -> QueryResult:
    _ensure_executable_exists(mysql_exe)
    sql = sql.strip()
    if not sql:
        return QueryResult(columns=[], rows=[], raw_stdout="")

    defaults_path = _write_temp_defaults_file(config._defaults_file_text())
    try:
        args = [
            mysql_exe,
            f"--defaults-extra-file={str(defaults_path)}",
            "--protocol=tcp",
            f"--connect-timeout={config.connect_timeout_sec}",
            "--batch",
            "--raw",
            "--silent",
            "--column-names",
            f"--execute={sql}",
        ]
        proc = _run_command(args, timeout_sec=timeout_sec)
        if proc.returncode != 0:
            raise MySQLCLIError("mysql 执行失败", command=args, stderr=proc.stderr.strip())

        stdout = proc.stdout or ""
        lines = [ln for ln in stdout.splitlines() if ln is not None]
        if not lines:
            return QueryResult(columns=[], rows=[], raw_stdout=stdout)

        columns = lines[0].split("\t")
        rows = [ln.split("\t") for ln in lines[1:]]
        return QueryResult(columns=columns, rows=rows, raw_stdout=stdout)
    finally:
        try:
            defaults_path.unlink(missing_ok=True)
        except OSError:
            pass


def test_connection(config: MySQLClientConfig, *, mysql_exe: str = "mysql") -> None:
    run_query(config, "SELECT 1;", mysql_exe=mysql_exe, timeout_sec=15)


def dump_sql(
    config: MySQLClientConfig,
    output_path: Path,
    *,
    tables: list[str] | None = None,
    mysqldump_exe: str = "mysqldump",
    timeout_sec: int = 300,
) -> None:
    _ensure_executable_exists(mysqldump_exe)
    if not config.database:
        raise ValueError("导出 .sql 需要填写 database")

    defaults_path = _write_temp_defaults_file(config._defaults_file_text())
    try:
        args: list[str] = [
            mysqldump_exe,
            f"--defaults-extra-file={str(defaults_path)}",
            "--protocol=tcp",
            "--single-transaction",
            "--quick",
            "--routines",
            "--events",
            config.database,
        ]
        if tables:
            args.extend(tables)

        output_path.parent.mkdir(parents=True, exist_ok=True)
        try:
            with output_path.open("wb") as f:
                proc = subprocess.run(
                    args,
                    stdout=f,
                    stderr=subprocess.PIPE,
                    timeout=timeout_sec,
                    check=False,
                )
        except subprocess.TimeoutExpired as e:
            raise MySQLCLIError(f"命令执行超时（>{timeout_sec}s）：{shlex.join(args)}") from e
        except OSError as e:
            raise MySQLCLIError(f"无法执行命令：{shlex.join(args)}") from e

        if proc.returncode != 0:
            stderr = (proc.stderr or b"").decode("utf-8", errors="replace").strip()
            raise MySQLCLIError("mysqldump 执行失败", command=args, stderr=stderr)
    finally:
        try:
            defaults_path.unlink(missing_ok=True)
        except OSError:
            pass
