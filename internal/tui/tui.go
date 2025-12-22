package tui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"mysql-client-gui/internal/db"
	"mysql-client-gui/internal/export"
)

type uiState struct {
	lastResult db.QueryResult
	hasResult  bool
}

func Run(logger *log.Logger) error {
	if logger == nil {
		return errors.New("logger 不能为空")
	}

	app := tview.NewApplication()
	state := &uiState{}

	cfg, err := db.DefaultConfigFromEnv()
	if err != nil {
		return err
	}

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]就绪")

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]快捷键：Tab/Shift+Tab 切换区域  Ctrl/Cmd+Enter 执行查询  Enter 确认  Esc 取消  Ctrl+C 退出")

	setStatus := func(s string) {
		status.SetText("[gray]" + s)
	}

	hostField := tview.NewInputField().SetLabel("Host: ").SetText(cfg.Host)
	portField := tview.NewInputField().SetLabel("Port: ").SetText(strconv.Itoa(cfg.Port))
	userField := tview.NewInputField().SetLabel("User: ").SetText(cfg.User)
	passField := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*')
	dbField := tview.NewInputField().SetLabel("Database: ").SetText(cfg.DB)

	sqlArea := tview.NewTextArea()
	sqlArea.SetText("SELECT 1;", true)
	sqlArea.SetBorder(true)
	sqlArea.SetTitle("SQL")

	resultTable := tview.NewTable()
	resultTable.SetFixed(1, 0)
	resultTable.SetBorders(false)
	resultTable.SetSelectable(true, true)
	resultTable.SetBorder(true)
	resultTable.SetTitle("结果")

	pages := tview.NewPages()

	showMessage := func(title, msg string) {
		modal := tview.NewModal().
			SetText(msg).
			AddButtons([]string{"确定"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("modal")
				app.SetFocus(sqlArea)
			})
		modal.SetTitle(title).SetBorder(true)
		pages.AddPage("modal", modal, true, true)
		app.SetFocus(modal)
	}

	promptPath := func(title, defaultPath string, onOK func(path string)) {
		input := tview.NewInputField().SetLabel("保存到: ").SetText(defaultPath)
		form := tview.NewForm().
			AddFormItem(input).
			AddButton("确定", func() {
				p := strings.TrimSpace(input.GetText())
				if p == "" {
					showMessage(title, "路径不能为空")
					return
				}
				pages.RemovePage("modal")
				onOK(p)
			}).
			AddButton("取消", func() {
				pages.RemovePage("modal")
				app.SetFocus(sqlArea)
			})
		form.SetBorder(true).SetTitle(title)
		pages.AddPage("modal", form, true, true)
		app.SetFocus(input)
	}

	buildConfig := func() (db.Config, error) {
		c := cfg
		c.Host = strings.TrimSpace(hostField.GetText())

		portText := strings.TrimSpace(portField.GetText())
		if portText == "" {
			c.Port = 3306
		} else {
			p, err := strconv.Atoi(portText)
			if err != nil {
				return db.Config{}, errors.New("port 必须是数字")
			}
			c.Port = p
		}

		c.User = strings.TrimSpace(userField.GetText())
		c.Pass = passField.GetText()
		c.DB = strings.TrimSpace(dbField.GetText())
		return c, c.Validate()
	}

	renderResult := func(res db.QueryResult) {
		resultTable.Clear()
		if len(res.Columns) == 0 {
			resultTable.SetCell(0, 0, tview.NewTableCell("(无结果)").SetTextColor(tview.Styles.SecondaryTextColor))
			return
		}

		for c, col := range res.Columns {
			resultTable.SetCell(0, c, tview.NewTableCell(col).SetSelectable(false).SetAttributes(tcell.AttrBold))
		}
		for r, row := range res.Rows {
			for c := range res.Columns {
				val := ""
				if c < len(row) {
					val = row[c]
				}
				resultTable.SetCell(r+1, c, tview.NewTableCell(val))
			}
		}
	}

	runQuery := func() {
		c, err := buildConfig()
		if err != nil {
			showMessage("执行查询", err.Error())
			return
		}
		sqlText := strings.TrimSpace(sqlArea.GetText())
		setStatus("执行中...")
		res, err := db.Run(c, sqlText)
		if err != nil {
			setStatus("执行查询：失败")
			showMessage("执行查询", err.Error())
			return
		}
		state.lastResult = res
		state.hasResult = len(res.Columns) > 0
		renderResult(res)
		setStatus(fmt.Sprintf("完成：%d 行", len(res.Rows)))
	}

	form := tview.NewForm().
		AddFormItem(hostField).
		AddFormItem(portField).
		AddFormItem(userField).
		AddFormItem(passField).
		AddFormItem(dbField).
		AddButton("测试连接", func() {
			c, err := buildConfig()
			if err != nil {
				showMessage("测试连接", err.Error())
				return
			}
			setStatus("测试连接中...")
			if err := db.TestConnection(c); err != nil {
				setStatus("测试连接：失败")
				showMessage("测试连接", err.Error())
				return
			}
			setStatus("连接正常")
			showMessage("测试连接", "连接成功")
		}).
		AddButton("执行查询", func() {
			runQuery()
		}).
		AddButton("导出 CSV(查询结果)", func() {
			if !state.hasResult || len(state.lastResult.Columns) == 0 {
				showMessage("导出 CSV", "暂无可导出的查询结果，请先执行查询。")
				return
			}
			promptPath("导出 CSV", "./result.csv", func(path string) {
				path = expandPath(path)
				setStatus("导出 CSV 中...")
				if err := export.WriteCSV(path, state.lastResult); err != nil {
					setStatus("导出 CSV：失败")
					showMessage("导出 CSV", err.Error())
					return
				}
				setStatus("已导出 CSV：" + path)
				showMessage("导出 CSV", "已保存："+path)
			})
		})

	form.SetBorder(true).SetTitle("连接")

	main := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 0, 3, true).
		AddItem(sqlArea, 0, 3, false).
		AddItem(resultTable, 0, 6, false).
		AddItem(help, 1, 0, false).
		AddItem(status, 1, 0, false)

	pages.AddPage("main", main, true, true)

	app.SetRoot(pages, true).SetFocus(sqlArea)
	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if pages.HasPage("modal") {
			return ev
		}
		if ev.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		if ev.Key() == tcell.KeyEnter && ev.Modifiers()&(tcell.ModCtrl|tcell.ModMeta) != 0 {
			runQuery()
			return nil
		}
		isShiftTab := ev.Key() == tcell.KeyBacktab || (ev.Key() == tcell.KeyTab && ev.Modifiers()&tcell.ModShift != 0)
		if ev.Key() == tcell.KeyTab || isShiftTab {
			if form.HasFocus() {
				formItem, button := form.GetFocusedItemIndex()
				items := form.GetFormItemCount()
				buttons := form.GetButtonCount()

				isFirst := false
				isLast := false
				if items > 0 {
					isFirst = formItem == 0 && button == -1
					if buttons > 0 {
						isLast = button == buttons-1
					} else {
						isLast = formItem == items-1
					}
				} else if buttons > 0 {
					isFirst = button == 0
					isLast = button == buttons-1
				}

				if isShiftTab && isFirst {
					app.SetFocus(resultTable)
					return nil
				}
				if !isShiftTab && isLast {
					app.SetFocus(sqlArea)
					return nil
				}
				return ev
			}

			order := []tview.Primitive{form, sqlArea, resultTable}
			current := app.GetFocus()
			next := 0
			for i, p := range order {
				if p == current {
					if isShiftTab {
						next = (i - 1 + len(order)) % len(order)
					} else {
						next = (i + 1) % len(order)
					}
					break
				}
			}
			if order[next] == form {
				items := form.GetFormItemCount()
				buttons := form.GetButtonCount()
				if isShiftTab {
					if buttons > 0 {
						form.SetFocus(items + buttons - 1)
					} else if items > 0 {
						form.SetFocus(items - 1)
					}
				} else if items > 0 {
					form.SetFocus(0)
				}
				app.SetFocus(form)
				return nil
			}
			app.SetFocus(order[next])
			return nil
		}
		return ev
	})

	// Ensure we have a valid CWD for relative exports.
	if _, err := os.Getwd(); err != nil {
		logger.Printf("WARN: Getwd failed: %v\n", err)
	}
	return app.Run()
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~/"))
		}
	}
	return p
}
