package internal

import (
    "github.com/lazybytez/jojo-discord-bot/components"
    "github.com/rs/zerolog/log"
    "os"
    "os/signal"
    "syscall"
)

// Bootstrap hadles the start of the application.
// It is responsible to execute the startup sequence
// and get the application up and running properly.
func Bootstrap() {
    startBot(Config.token)

    components.RegisterComponents(discord)

    waitForTerminate()
}

// waitForTerminate blocks the console and waits
// for a termination signal.
//
// When a sigterm is received, the application is stopped
// gracefully. This means all open connections or used resources are
// freed/closed before exit.
func waitForTerminate() {
    log.Info().Msg("Bot is running.  Press CTRL-C to exit.")

    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-c

    ExitGracefully("Bot has been terminated gracefully!")
}
