# mysql-client-gui

封装 Linux 下 `mysql`/`mysqldump` 命令行客户端调用，提供一个轻量 GUI：

- 连接信息配置（host/port/user/password/database）
- SQL 查询执行与结果表格展示
- 将查询结果导出为 CSV
- 可选：使用 `mysqldump` 导出为 `.sql`

## 环境要求

- Linux
- Python 3.10+
- 系统已安装 MySQL CLI：`mysql`（以及可选 `mysqldump`）
- Tkinter（多数发行版自带；若缺失，安装 `python3-tk`）

## 运行

方式 1：直接运行

```bash
python3 -m mysql_client_gui
```

方式 2：安装为本地命令（可选）

```bash
python3 -m pip install -e .
mysql-client-gui
```

## 打包（CentOS 7，目标机无需 Python）

你可以在“有 Python 的打包机”上把程序打成 Linux 可执行目录（onedir），再拷贝到目标 CentOS 7 机器运行。

### 1) 在 CentOS 7 打包（推荐在 CentOS 7 x86_64 上执行）

要求：

- Python 3.10+（仅打包机需要）
- `mysql` CLI（目标机也需要）
- Tkinter（打包机需要；一般包名为 `python3-tkinter` / `python3-tk`，随发行版不同）

执行：

```bash
./scripts/build-linux-onedir.sh
```

输出：

- `dist/mysql-client-gui/`（拷贝这个目录到目标机）

### 1.5) 用 Docker 在 macOS/Windows 上打包（Ubuntu 22.04 环境）

如果本机不是 Linux（例如 macOS），可以用容器构建（默认 Ubuntu 22.04）。
注意：在 Ubuntu 22.04 里构建出来的二进制，通常要求目标机 glibc 版本不低于 Ubuntu 22.04（不保证兼容 CentOS 7）。

构建镜像：

```bash
docker build -f Dockerfile.linux-build -t mysql-client-gui:linux-build .
```

运行构建（把产物输出到本机 `./dist`）：

```bash
mkdir -p dist
docker run --rm -v "$PWD/dist:/out" mysql-client-gui:linux-build
```

Apple Silicon（M1/M2/M3）请用 x86_64 平台构建/运行：

```bash
docker buildx build --platform=linux/amd64 -f Dockerfile.linux-build -t mysql-client-gui:linux-build .
docker run --rm --platform=linux/amd64 -v "$PWD/dist:/out" mysql-client-gui:linux-build
```

### 2) 在目标机运行

把 `dist/mysql-client-gui/` 整个目录拷贝到目标机后执行：

```bash
./mysql-client-gui
```

目标机注意事项：

- 需要可用的图形环境（X11/Wayland）；否则 Tk 窗口无法打开
- 需要 `mysql` 在 `PATH` 中（本项目是封装系统 mysql 客户端，不会把 mysql 一起打包进来）

## 使用提示

- 密码不会写入磁盘持久化；程序会为每次调用生成临时 `--defaults-extra-file`，执行完即删除。
- 查询结果展示会把第一行当做列名（`mysql --batch` 输出）。
