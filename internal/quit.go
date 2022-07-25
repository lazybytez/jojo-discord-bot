package internal

import (
    "github.com/rs/zerolog/log"
    "os"
)

func ExitGracefull(reason string) {
    log.Info().Msg(reason)
    os.Exit(0)
}

func ExitFatal(reason string) {
    log.Fatal().Msg(reason)
}
