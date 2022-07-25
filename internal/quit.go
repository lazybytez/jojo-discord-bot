package internal

import (
    "github.com/rs/zerolog/log"
    "os"
)

// ExitGracefully shutdowns the bot gracefully.
// Using this functions results in the process exiting with
// exit code 0.
//
// The function tries to ensure that all allocated resources like
// connections or locks are freed/closed correctly.
func ExitGracefully(reason string) {
    releaseResources()

    log.Info().Msg(reason)
    os.Exit(0)
}

// ExitFatal shutdowns the application ungracefully with a
// non-zero exit code.
// The routines to properly stop the application are not applied on
// a fatal exit. The function should be only called when the application cannot
// recover, which is typically the case when core connections cannot be established
// or an initialization routine fails.
func ExitFatal(reason string) {
    log.Fatal().Msg(reason)
}

// releaseResources ensures that all allocated resources, locks
// and connections are freed before the application terminates
// gracefully.
func releaseResources() {
    if nil != discord {
        _ = discord.Close()
    }
}
