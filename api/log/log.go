package log

import "github.com/rs/zerolog/log"

// ComponentLogPrefix is used by the logger to prefix
// log messages with the component name.
const ComponentLogPrefix = "component"

// Debug logs a message using the applications logger.
// The log level of the logged message is debug.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Debug(componentName string, format string, v ...interface{}) {
    if len(v) == 0 {
        log.Debug().Str(ComponentLogPrefix, componentName).Msg(format)

        return
    }

    log.Debug().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}

// Info logs a message using the applications logger.
// The log level of the logged message is info.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Info(componentName string, format string, v ...interface{}) {
    if len(v) == 0 {
        log.Info().Str(ComponentLogPrefix, componentName).Msg(format)

        return
    }

    log.Info().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}

// Warn logs a message using the applications logger.
// The log level of the logged message is warning.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Warn(componentName string, format string, v ...interface{}) {
    if len(v) == 0 {
        log.Warn().Str(ComponentLogPrefix, componentName).Msg(format)

        return
    }

    log.Warn().Str(ComponentLogPrefix, componentName).Msgf(format, v)
}

// Err logs a message using the applications logger.
// The log level of the logged message is error.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Err(componentName string, err error, format string, v ...interface{}) {
    if len(v) == 0 {
        log.Err(err).Str(ComponentLogPrefix, componentName).Msg(format)

        return
    }

    log.Error().Err(err).Str(ComponentLogPrefix, componentName).Msgf(format, v)
}
