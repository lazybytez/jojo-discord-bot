package internal

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog/log"
    "os"
    "path"
)

const envFile = ".env"

const token = "TOKEN"

// JojoBotConfig represents the entire environment variable based configuration
// of the application. Globally available values are public,
// sensitive values are private and can only being read
// from within the current package.
//
// When new configuration options get added to the application, they should be appended
// to the structure.
type JojoBotConfig struct {
    token string
}

// Config holds the currently loaded configuration
// of the application.
//
// It should be used to retrieve the set environment variables
// instead of directly calling os.Getenv.
//
// The variable is initialized during the init function.
var Config JojoBotConfig

// Initialize environment with local .env file
// This will load the environment variables defined in the specified
// env file and merge them into os.Environ.
func init() {
    pwd, err := os.Getwd()
    if nil != err {
        ExitFatal("Failed to get current working directory to load env file from!")
    }

    envFilePath := path.Join(pwd, envFile)
    log.Info().Msgf("Trying to load env file from \"%v\"...", envFilePath)

    err = godotenv.Load(envFilePath)
    if nil != err {
        ExitFatal(fmt.Sprintf("Missing env file at %v", envFilePath))
    }
    log.Info().Msgf("Sucessfully loaded env file from \"%v\"!", envFilePath)

    Config = JojoBotConfig{
        token: os.Getenv(token),
    }

    cleanUpSensitiveValues()
}

// cleanUpSensitiveValues throws out sensitive values from the
// environment configuration to prevent other packages from using those
// without reloading them. This reduces the attack surface from accidentally
// printing out the env to a user of the application.
//
// Note that variables loaded using godotenv.Load could be loaded again at any time
func cleanUpSensitiveValues() {
    err := os.Setenv(token, "")

    if nil != err {
        ExitFatal("Failed to cleanup sensitive environment values!")
    }
}
