package main

import log "github.com/sirupsen/logrus"

// CustomLevelColors is a custom hook to set colors for specific log levels
type CustomLevelColors struct {
	InfoColor  string
	ErrorColor string
	WarnColor  string
}

// Fire sets colors for Info and Error log levels
func (hook *CustomLevelColors) Fire(entry *log.Entry) error {
	switch entry.Level {
	case log.InfoLevel:
		entry.Message = hook.InfoColor + entry.Message + "[0m" // Reset color after the message
	case log.ErrorLevel:
		entry.Message = hook.ErrorColor + entry.Message + "[0m" // Reset color after the message
	}
	return nil
}

// Levels returns the log levels for the hook
func (hook *CustomLevelColors) Levels() []log.Level {
	return []log.Level{log.InfoLevel, log.ErrorLevel, log.WarnLevel}
}

func init() {
	logger := log.New()
	// Optionally, set the log level
	logger.SetLevel(log.DebugLevel)
	// Create a custom formatter with force colors enabled
	logger.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	// Add a custom hook to set colors for specific log levels
	logger.AddHook(&CustomLevelColors{
		InfoColor:  "[34m", // Blue for Info
		ErrorColor: "[31m", // Red for Error
		WarnColor:  "[33m", // Yellow for Warn
	})
}
