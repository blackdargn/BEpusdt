package cmd

import (
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/conf"
)

type RuntimeSettings struct {
	Listen   string
	Log      string
	SQLite   string
	MySQL    string
	Postgres string
}

func resolveSettings(c *cli.Command) RuntimeSettings {
	cfg := conf.GetFileConfig()

	return RuntimeSettings{
		Listen:   resolveString(c, "listen", "LISTEN", cfgValue(cfg, func(c *conf.FileConfig) string { return c.Listen }), ":8080"),
		Log:      resolveLog(c, cfg),
		SQLite:   resolveString(c, "sqlite", "SQLITE", cfgValue(cfg, func(c *conf.FileConfig) string { return c.SQLite }), "/var/lib/bepusdt/sqlite.db"),
		MySQL:    resolveString(c, "mysql", "MYSQL_DSN", cfgValue(cfg, func(c *conf.FileConfig) string { return c.MySQLDSN }), ""),
		Postgres: resolveString(c, "postgres", "POSTGRESQL_DSN", cfgValue(cfg, func(c *conf.FileConfig) string { return c.PostgreSQLDSN }), ""),
	}
}

func resolveLog(c *cli.Command, cfg *conf.FileConfig) string {
	if cliFlagPassed(c, "log") {
		return c.String("log")
	}

	if value := strings.TrimSpace(cfgValue(cfg, func(c *conf.FileConfig) string { return c.Log })); value != "" {
		return value
	}

	return conf.DefaultLogDir()
}

func resolveString(c *cli.Command, flagName, envName, configVal, defaultVal string) string {
	if cliFlagPassed(c, flagName) {
		return c.String(flagName)
	}

	if value := strings.TrimSpace(os.Getenv(envName)); value != "" {
		return value
	}

	if value := strings.TrimSpace(configVal); value != "" {
		return value
	}

	return defaultVal
}

func cfgValue(cfg *conf.FileConfig, getter func(*conf.FileConfig) string) string {
	if cfg == nil {
		return ""
	}

	return getter(cfg)
}

func cliFlagPassed(c *cli.Command, name string) bool {
	for _, arg := range commandArgs(c.Name) {
		if arg == "--"+name || strings.HasPrefix(arg, "--"+name+"=") {
			return true
		}
	}

	return false
}

func commandArgs(subcommand string) []string {
	for i, arg := range os.Args {
		if arg != subcommand {
			continue
		}

		if i+1 >= len(os.Args) {
			return nil
		}

		return os.Args[i+1:]
	}

	return nil
}
