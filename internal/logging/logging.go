package logging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

// L is the global logrus logger used across the application.
var L *logrus.Logger

func newLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stderr)
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		FullTimestamp:    false,
		ForceColors:      term.IsTerminal(int(os.Stderr.Fd())),
		PadLevelText:     true,
	})
	return l
}

// Setup configures the global logger with the provided level.
// Supported levels: debug, info, warn, error.
func Setup(level string) {
	L = newLogger()
	lvl := parseLevel(level)
	L.SetLevel(lvl)
	if tf, ok := L.Formatter.(*logrus.TextFormatter); ok {
		tf.ForceColors = term.IsTerminal(int(os.Stderr.Fd()))
		tf.DisableTimestamp = true
	}
}

func parseLevel(level string) logrus.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

// Convert alternating key/value args to logrus.Fields
func toFields(args ...any) logrus.Fields {
	if len(args) == 0 {
		return nil
	}
	fields := logrus.Fields{}
	for i := 0; i+1 < len(args); i += 2 {
		key := fmt.Sprint(args[i])
		fields[key] = args[i+1]
	}
	if len(args)%2 == 1 {
		fields["arg"] = args[len(args)-1]
	}
	return fields
}

// Convenience helpers. Prefer using these from callers.
func Debug(msg string, args ...any) { L.WithFields(toFields(args...)).Debug(msg) }
func Info(msg string, args ...any)  { L.WithFields(toFields(args...)).Info(msg) }
func Warn(msg string, args ...any)  { L.WithFields(toFields(args...)).Warn(msg) }
func Error(msg string, args ...any) { L.WithFields(toFields(args...)).Error(msg) }

// Formatted logging helpers
func Debugf(format string, args ...any) { L.Debugf(format, args...) }
func Infof(format string, args ...any)  { L.Infof(format, args...) }
func Warnf(format string, args ...any)  { L.Warnf(format, args...) }
func Errorf(format string, args ...any) { L.Errorf(format, args...) }

// Context-aware helpers (ctx is unused in logrus core, kept for API convenience)
func DebugContext(ctx context.Context, msg string, args ...any) { Debug(msg, args...) }
func InfoContext(ctx context.Context, msg string, args ...any)  { Info(msg, args...) }
func WarnContext(ctx context.Context, msg string, args ...any)  { Warn(msg, args...) }
func ErrorContext(ctx context.Context, msg string, args ...any) { Error(msg, args...) }
