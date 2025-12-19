package db

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Host string
	Port int
	User string
	Pass string
	DB   string

	ConnectTimeout time.Duration
	QueryTimeout   time.Duration

	Charset string
	Loc     string
}

func DefaultConfig() Config {
	return Config{
		Host:           "127.0.0.1",
		Port:           3306,
		User:           "root",
		Pass:           "",
		DB:             "",
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   60 * time.Second,
		Charset:        "utf8mb4",
		Loc:            "Local",
	}
}

func DefaultConfigFromEnv() (Config, error) {
	return ApplyEnv(DefaultConfig())
}

func ApplyEnv(c Config) (Config, error) {
	if v := strings.TrimSpace(os.Getenv("MYSQL_HOST")); v != "" {
		c.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("MYSQL_PORT")); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, errors.New("MYSQL_PORT 必须是数字")
		}
		c.Port = p
	}
	if v := strings.TrimSpace(os.Getenv("MYSQL_USER")); v != "" {
		c.User = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		c.Pass = v
	}
	if v := strings.TrimSpace(os.Getenv("MYSQL_DB")); v != "" {
		c.DB = v
	} else if v := strings.TrimSpace(os.Getenv("MYSQL_DATABASE")); v != "" {
		c.DB = v
	}
	return c, nil
}

func (c Config) Validate() error {
	host := strings.TrimSpace(c.Host)
	user := strings.TrimSpace(c.User)
	if host == "" {
		return errors.New("host 不能为空")
	}
	if user == "" {
		return errors.New("user 不能为空")
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port 范围应为 1..65535")
	}
	if strings.Contains(host, "/") {
		return errors.New("host 格式不合法")
	}
	return nil
}

func (c Config) Addr() string {
	return net.JoinHostPort(strings.TrimSpace(c.Host), strconv.Itoa(c.Port))
}

func (c Config) DSN() string {
	dbName := strings.TrimSpace(c.DB)
	user := strings.TrimSpace(c.User)
	pass := c.Pass

	mc := mysql.NewConfig()
	mc.User = user
	mc.Passwd = pass
	mc.Net = "tcp"
	mc.Addr = c.Addr()
	mc.DBName = dbName // can be empty
	mc.Timeout = c.ConnectTimeout
	mc.ReadTimeout = c.QueryTimeout
	mc.WriteTimeout = c.QueryTimeout
	mc.MultiStatements = true
	mc.ParseTime = true
	mc.Params = map[string]string{
		"charset": c.Charset,
	}
	if loc, err := time.LoadLocation(c.Loc); err == nil {
		mc.Loc = loc
	}
	return mc.FormatDSN()
}
