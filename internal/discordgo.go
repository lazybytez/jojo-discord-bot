package internal

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session

// startBot initializes a discordgo.Session using
// the provided token.
//
// If initializing the discordgo.Session fails,
// the bot will exit with a fatal error.
func startBot(token string) {
    var err error
    discord, err = discordgo.New("Bot " + token)
    if nil != err {
        ExitFatal(fmt.Sprintf("Failed to initialize Discord connection, error was: %v!", err.Error()))
    }

    fmt.Println(discord.Token)
}

// stopBot tries to stop the bot.
// The bot is stopped by closing the discordgo.Session.
//
// If the bot has not been initialized until this point,
// the close function of the discordgo.Session won't be called.
//
// If closing the session throws an error, the error is ignored.
func stopBot() {
    if nil != discord {
        _ = discord.Close()
    }
}
