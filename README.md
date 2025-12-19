# mysql-client-gui

一个轻量 MySQL 客户端（Go 实现，终端 UI：`tview`），用于查询与导出：

- 连接信息配置（host/port/user/password/database）
- SQL 查询执行与结果表格展示
- 将查询结果导出为 CSV

## 环境要求

- Linux
- 运行：可用的终端（无需图形环境）
- 构建：Go 1.22+

## 运行

方式 1：从源码运行

```bash
go run ./cmd/mysql-client-gui
```

或使用 Makefile：

```bash
make run
```

方式 2：构建后运行（推荐）

```bash
make build
./bin/mysql-client-gui
```

## 打包（Linux，目标机无需 Go）

你可以在“有 Go 的打包机”上构建二进制，再拷贝到目标 Linux 机器运行。

### 1) 在 Linux 上打包（推荐在较老 Linux 上执行以获得更高兼容性）

执行：

```bash
make build
```

如需构建不同架构，可指定：`make build GOARCH=arm64`（默认 `linux/amd64`）。

输出：

- `bin/mysql-client-gui`

### 2) 在目标机运行

把 `bin/mysql-client-gui` 拷贝到目标机后执行：

```bash
./mysql-client-gui
```

### 命令行帮助

```bash
./mysql-client-gui --help
```

### 环境变量（可选）

可通过环境变量预设连接信息，启动后会自动填入输入框：

- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DB`（或 `MYSQL_DATABASE`）

## 使用提示

### 基本用法

1) 填写连接信息（Host/Port/User/Password/Database）。  
2) 点击“测试连接”确认配置正确。  
3) 在 SQL 区输入语句，点击“执行查询”。结果会显示在表格中。  
4) 点击“导出 CSV(查询结果)”并选择保存路径（默认 `./result.csv`，支持 `~/`）。  

### 快捷键

- `Tab`/`Shift+Tab`：切换输入区域  
- `Enter`：确认/执行  
- `Esc`：取消弹窗  
- `Ctrl+C`：退出程序  

### 安全提示

- 密码不会写入磁盘持久化（仅在内存中使用）。
