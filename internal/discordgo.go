package internal

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session

func startBot(token string) {
    var err error
    discord, err = discordgo.New(token)
    if nil != err {
        ExitFatal(fmt.Sprintf("Failed to initialize Discord connection, error was: %v!", err.Error()))
    }

    fmt.Println(discord.Token)
}
