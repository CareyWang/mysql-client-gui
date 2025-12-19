package export

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"

	"mysql-client-gui/internal/db"
)

func WriteCSV(path string, res db.QueryResult) error {
	if len(res.Columns) == 0 {
		return errors.New("暂无可导出的查询结果，请先执行查询。")
	}
	if path == "" {
		return errors.New("路径不能为空")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(res.Columns); err != nil {
		return err
	}
	for _, r := range res.Rows {
		if err := w.Write(r); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
