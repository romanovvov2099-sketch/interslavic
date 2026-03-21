package logging

import (
	"fmt"
	"log/slog"
)

type Prefixer interface {
	GetPrefix() string
}

type BaseLogger struct {
	config *LoggingConfig

	logDirPath string

	version     string
	moduleIdent string // VERSION:${version}
}

func NewBaseLogger(version string, dir string, c *LoggingConfig) *BaseLogger {
	return &BaseLogger{
		config: c,

		logDirPath: dir,

		version:     version,
		moduleIdent: fmt.Sprintf("%s:%s", baseModuleName, version),
	}
}

func InvokeBaseLogger(bl *BaseLogger) {
	bl.InvokeLogging()
}

func (bl *BaseLogger) GetPrefix() string {
	return bl.moduleIdent
}

type ModuleLogger struct {
	*slog.Logger

	moduleName  string
	moduleValue string
	moduleIdent string // moduleName:moduleValue
}

func NewModuleLogger(n string, v string, p Prefixer) *ModuleLogger {
	identifier := fmt.Sprintf("%s:%s", n, v)
	if p != nil {
		identifier = fmt.Sprintf("%s | %s", p.GetPrefix(), identifier)
	}

	return &ModuleLogger{
		Logger: slog.Default(),

		moduleIdent: identifier,
		moduleName:  n,
		moduleValue: v,
	}
}

func (l *ModuleLogger) GetPrefix() string {
	return l.moduleIdent
}

func (l *ModuleLogger) Debug(msg string, args ...any) {
	baseMsg := fmt.Sprintf("%s - %s", l.moduleIdent, msg)

	l.Logger.Debug(baseMsg, args...)
}

func (l *ModuleLogger) Info(msg string, args ...any) {
	baseMsg := fmt.Sprintf("%s - %s", l.moduleIdent, msg)

	l.Logger.Info(baseMsg, args...)
}

func (l *ModuleLogger) Warn(msg string, args ...any) {
	baseMsg := fmt.Sprintf("%s - %s", l.moduleIdent, msg)

	l.Logger.Warn(baseMsg, args...)
}

func (l *ModuleLogger) Error(msg string, args ...any) {
	baseMsg := fmt.Sprintf("%s - %s", l.moduleIdent, msg)

	l.Logger.Error(baseMsg, args...)
}

func (l *ModuleLogger) Critical(msg string, args ...any) {
	baseMsg := fmt.Sprintf("%s - %s", l.moduleIdent, msg)

	l.Logger.Error(baseMsg, args...)
	panic(baseMsg)
}

// Notify using [MsgMainKeyword = "DRIVER"]
func Notify(keyword string) {
	msg := fmt.Sprintf("------ %s %s ------", MsgMainKeyword, keyword)

	slog.Info(msg)
}

func NotifyWithOneMsg(keymsg string) {
	msg := fmt.Sprintf("------ %s ------", keymsg)

	slog.Info(msg)
}