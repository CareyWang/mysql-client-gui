package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mysql-client-gui/internal/tui"
)

func main() {
	showHelp := false
	flag.BoolVar(&showHelp, "help", false, "显示帮助")
	flag.BoolVar(&showHelp, "h", false, "显示帮助")
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, "mysql-client-gui - 轻量 MySQL 客户端（TUI）")
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "用法：")
		fmt.Fprintln(os.Stdout, "  mysql-client-gui [--help|-h]")
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "环境变量（可选）：")
		fmt.Fprintln(os.Stdout, "  MYSQL_HOST       MySQL Host")
		fmt.Fprintln(os.Stdout, "  MYSQL_PORT       MySQL Port")
		fmt.Fprintln(os.Stdout, "  MYSQL_USER       MySQL User")
		fmt.Fprintln(os.Stdout, "  MYSQL_PASSWORD   MySQL Password")
		fmt.Fprintln(os.Stdout, "  MYSQL_DB         MySQL Database")
		fmt.Fprintln(os.Stdout, "  MYSQL_DATABASE   MySQL Database（兼容别名）")
	}
	flag.Parse()
	if showHelp {
		flag.Usage()
		return
	}

	logger := log.New(os.Stderr, "", 0)
	if err := tui.Run(logger); err != nil {
		logger.Printf("错误：%v\n", err)
		os.Exit(1)
	}
}
