package applogger

import "github.com/hashicorp/go-hclog"

func NewLogger(applicationName string) hclog.Logger {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  applicationName,
		Level: hclog.LevelFromString("DEBUG"),
	})
	return appLogger
}
