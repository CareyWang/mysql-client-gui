package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type QueryResult struct {
	Columns []string
	Rows    [][]string
}

func Open(cfg Config) (*sql.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	return db, nil
}

func TestConnection(cfg Config) error {
	db, err := Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()
	return db.PingContext(ctx)
}

func Run(cfg Config, sqlText string) (QueryResult, error) {
	if err := cfg.Validate(); err != nil {
		return QueryResult{}, err
	}
	sqlText = strings.TrimSpace(sqlText)
	if sqlText == "" {
		return QueryResult{Columns: nil, Rows: nil}, nil
	}

	db, err := Open(cfg)
	if err != nil {
		return QueryResult{}, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.QueryTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, sqlText)
	if err != nil {
		// Not all statements return rows; fall back to Exec.
		if _, execErr := db.ExecContext(ctx, sqlText); execErr == nil {
			return QueryResult{Columns: nil, Rows: nil}, nil
		}
		return QueryResult{}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}

	colTypes, _ := rows.ColumnTypes()
	_ = colTypes

	dest := make([]any, len(cols))
	raw := make([]any, len(cols))
	for i := range dest {
		dest[i] = &raw[i]
	}

	var outRows [][]string
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return QueryResult{}, err
		}
		row := make([]string, 0, len(cols))
		for _, v := range raw {
			row = append(row, formatCell(v))
		}
		outRows = append(outRows, row)
	}
	if err := rows.Err(); err != nil {
		return QueryResult{}, err
	}

	return QueryResult{Columns: cols, Rows: outRows}, nil
}

func formatCell(v any) string {
	if v == nil {
		return "NULL"
	}
	switch x := v.(type) {
	case []byte:
		return string(x)
	default:
		return fmt.Sprint(v)
	}
}
