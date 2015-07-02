package newrelic

import (
	l "log"
	"os"
)

// LoggingLevel enumerates package log levels
type LoggingLevel int

const (
	// LogAll logs verbosely
	LogAll LoggingLevel = 0
	// LogDebug logs debug or above
	LogDebug LoggingLevel = 20
	// LogInfo logs informational or above
	LogInfo LoggingLevel = 30
	// LogError only logs errors
	LogError LoggingLevel = 50
	// LogNone makes the package silent
	LogNone LoggingLevel = 100
)

// Logger is the logger used by this package. Set to a custom logger if needed.
var Logger = l.New(os.Stderr, "newrelic", l.LstdFlags)

// LogLevel can be set globally to tune logging levels
var LogLevel = LogError

func log(level LoggingLevel, format string, a ...interface{}) {
	if level <= LogLevel {
		Logger.Printf(format, a...)
	}
}
