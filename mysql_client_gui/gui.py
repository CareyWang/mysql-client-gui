from __future__ import annotations

import csv
from pathlib import Path
from tkinter import END, BOTH, BOTTOM, LEFT, RIGHT, TOP, X, Y, StringVar, Tk, ttk, filedialog, messagebox
from tkinter.scrolledtext import ScrolledText

from mysql_client_gui.mysql_cli import MySQLCLIError, MySQLClientConfig, QueryResult, dump_sql, run_query, test_connection


def _parse_port(port_text: str) -> int:
    port_text = port_text.strip()
    if not port_text:
        return 3306
    port = int(port_text)
    if not (1 <= port <= 65535):
        raise ValueError("port 范围应为 1..65535")
    return port


def _build_config(host: str, port: str, user: str, password: str, database: str) -> MySQLClientConfig:
    host = host.strip()
    user = user.strip()
    database = database.strip()
    if not host:
        raise ValueError("host 不能为空")
    if not user:
        raise ValueError("user 不能为空")
    db = database or None
    return MySQLClientConfig(
        host=host,
        port=_parse_port(port),
        user=user,
        password=password,
        database=db,
    )


class App:
    def __init__(self, root: Tk) -> None:
        self.root = root
        self.root.title("MySQL Client GUI (mysql CLI wrapper)")
        self.root.geometry("1100x700")

        self.host_var = StringVar(value="127.0.0.1")
        self.port_var = StringVar(value="3306")
        self.user_var = StringVar(value="root")
        self.pass_var = StringVar(value="")
        self.db_var = StringVar(value="")
        self.dump_tables_var = StringVar(value="")
        self.status_var = StringVar(value="就绪")

        self._last_result: QueryResult | None = None

        self._build_ui()

    def _build_ui(self) -> None:
        conn = ttk.LabelFrame(self.root, text="连接")
        conn.pack(side=TOP, fill=X, padx=10, pady=8)

        def add_row(label: str, var: StringVar, *, show: str | None = None) -> ttk.Entry:
            row = ttk.Frame(conn)
            row.pack(side=TOP, fill=X, padx=8, pady=4)
            ttk.Label(row, text=label, width=12).pack(side=LEFT)
            entry = ttk.Entry(row, textvariable=var, show=show) if show else ttk.Entry(row, textvariable=var)
            entry.pack(side=LEFT, fill=X, expand=True)
            return entry

        add_row("Host", self.host_var)
        add_row("Port", self.port_var)
        add_row("User", self.user_var)
        add_row("Password", self.pass_var, show="*")
        add_row("Database", self.db_var)

        buttons = ttk.Frame(conn)
        buttons.pack(side=TOP, fill=X, padx=8, pady=6)
        ttk.Button(buttons, text="测试连接", command=self.on_test).pack(side=LEFT)
        ttk.Button(buttons, text="执行查询", command=self.on_run).pack(side=LEFT, padx=8)
        ttk.Button(buttons, text="导出 CSV(查询结果)", command=self.on_export_csv).pack(side=LEFT)
        ttk.Button(buttons, text="导出 SQL(mysqldump)", command=self.on_dump_sql).pack(side=LEFT, padx=8)

        query_frame = ttk.LabelFrame(self.root, text="SQL")
        query_frame.pack(side=TOP, fill=BOTH, expand=False, padx=10, pady=8)
        self.sql_text = ScrolledText(query_frame, height=8, wrap="none")
        self.sql_text.pack(fill=BOTH, expand=True, padx=8, pady=6)
        self.sql_text.insert(END, "SELECT 1;")

        dump_frame = ttk.LabelFrame(self.root, text="mysqldump（可选）")
        dump_frame.pack(side=TOP, fill=X, padx=10, pady=8)
        row = ttk.Frame(dump_frame)
        row.pack(side=TOP, fill=X, padx=8, pady=4)
        ttk.Label(row, text="Tables(可选，用逗号分隔)", width=24).pack(side=LEFT)
        ttk.Entry(row, textvariable=self.dump_tables_var).pack(side=LEFT, fill=X, expand=True)

        result_frame = ttk.LabelFrame(self.root, text="结果")
        result_frame.pack(side=TOP, fill=BOTH, expand=True, padx=10, pady=8)

        self.tree = ttk.Treeview(result_frame, columns=(), show="headings")
        vsb = ttk.Scrollbar(result_frame, orient="vertical", command=self.tree.yview)
        hsb = ttk.Scrollbar(result_frame, orient="horizontal", command=self.tree.xview)
        self.tree.configure(yscrollcommand=vsb.set, xscrollcommand=hsb.set)
        vsb.pack(side=RIGHT, fill=Y)
        hsb.pack(side=BOTTOM, fill=X)
        self.tree.pack(side=LEFT, fill=BOTH, expand=True)

        status = ttk.Frame(self.root)
        status.pack(side=TOP, fill=X, padx=10, pady=6)
        ttk.Label(status, textvariable=self.status_var).pack(side=LEFT)

    def _set_status(self, text: str) -> None:
        self.status_var.set(text)
        self.root.update_idletasks()

    def _get_config(self) -> MySQLClientConfig:
        return _build_config(
            self.host_var.get(),
            self.port_var.get(),
            self.user_var.get(),
            self.pass_var.get(),
            self.db_var.get(),
        )

    def _handle_error(self, title: str, err: Exception) -> None:
        if isinstance(err, MySQLCLIError):
            detail = err.stderr or str(err)
        else:
            detail = str(err)
        messagebox.showerror(title, detail)
        self._set_status(f"{title}：失败")

    def on_test(self) -> None:
        try:
            cfg = self._get_config()
            self._set_status("测试连接中...")
            test_connection(cfg)
            self._set_status("连接正常")
            messagebox.showinfo("测试连接", "连接成功")
        except Exception as e:
            self._handle_error("测试连接", e)

    def on_run(self) -> None:
        try:
            cfg = self._get_config()
            sql = self.sql_text.get("1.0", END).strip()
            self._set_status("执行中...")
            result = run_query(cfg, sql)
            self._last_result = result
            self._render_result(result)
            self._set_status(f"完成：{len(result.rows)} 行")
        except Exception as e:
            self._handle_error("执行查询", e)

    def _render_result(self, result: QueryResult) -> None:
        self.tree.delete(*self.tree.get_children())
        self.tree["columns"] = tuple(result.columns)

        for col in result.columns:
            self.tree.heading(col, text=col)
            self.tree.column(col, width=160, stretch=True, anchor="w")

        for row in result.rows:
            values = row + [""] * max(0, len(result.columns) - len(row))
            self.tree.insert("", END, values=values[: len(result.columns)])

    def on_export_csv(self) -> None:
        if not self._last_result or not self._last_result.columns:
            messagebox.showwarning("导出 CSV", "暂无可导出的查询结果，请先执行查询。")
            return

        path = filedialog.asksaveasfilename(
            title="保存 CSV",
            defaultextension=".csv",
            filetypes=[("CSV", "*.csv"), ("All Files", "*.*")],
        )
        if not path:
            return

        try:
            out = Path(path)
            out.parent.mkdir(parents=True, exist_ok=True)
            with out.open("w", newline="", encoding="utf-8") as f:
                w = csv.writer(f)
                w.writerow(self._last_result.columns)
                w.writerows(self._last_result.rows)
            self._set_status(f"已导出 CSV：{out}")
            messagebox.showinfo("导出 CSV", f"已保存：{out}")
        except Exception as e:
            self._handle_error("导出 CSV", e)

    def on_dump_sql(self) -> None:
        try:
            cfg = self._get_config()
            if not cfg.database:
                raise ValueError("导出 SQL 需要填写 database")

            path = filedialog.asksaveasfilename(
                title="保存 SQL",
                defaultextension=".sql",
                filetypes=[("SQL", "*.sql"), ("All Files", "*.*")],
            )
            if not path:
                return

            tables_text = self.dump_tables_var.get().strip()
            tables = [t.strip() for t in tables_text.split(",") if t.strip()] if tables_text else None
            self._set_status("mysqldump 导出中...")
            dump_sql(cfg, Path(path), tables=tables)
            self._set_status("mysqldump 导出完成")
            messagebox.showinfo("导出 SQL", f"已保存：{path}")
        except Exception as e:
            self._handle_error("导出 SQL", e)


def main() -> None:
    root = Tk()
    style = ttk.Style(root)
    try:
        style.theme_use("clam")
    except Exception:
        pass
    App(root)
    root.mainloop()
