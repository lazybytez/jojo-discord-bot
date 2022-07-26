package log

import "github.com/rs/zerolog/log"

// ComponentLogPrefix is used by the logger to prefix
// log messages with the component name.
const ComponentLogPrefix = "component"

// LogDebug logs a message using the applications logger.
// The log level of the logged message is debug.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func LogDebug(componentName string, format string, v ...interface{}) {
    log.Debug().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}

// LogInfo logs a message using the applications logger.
// The log level of the logged message is info.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func LogInfo(componentName string, format string, v ...interface{}) {
    log.Info().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}

// LogWarn logs a message using the applications logger.
// The log level of the logged message is warning.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func LogWarn(componentName string, format string, v ...interface{}) {
    log.Warn().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}
