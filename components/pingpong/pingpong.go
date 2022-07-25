package pingpong

import (
    "github.com/bwmarrin/discordgo"
    "github.com/lazybytez/jojo-discord-bot/internal/component"
    "github.com/rs/zerolog/log"
)

// Component struct that acts as the type for the
// Ping-Pong component of the application.
type Component struct {
    component.Component
}

// ComponentInstance is the instance of the ping pong component.
// Can be used to register the component or get information about it.
var ComponentInstance = &Component{
    component.Component{
        Name:         "Ping Pong Component",
        Description:  "This Component plays pingpong with you and returns Laytency (maybe)",
        DmPermission: true,
    },
}

// LoadComponent loads the Ping-Pong Component
func (c *Component) LoadComponent(discord *discordgo.Session) (*component.Component, error) {
    // Register the messageCreate func as a callback for MessageCreate events.
    discord.AddHandler(onMessageCreate)

    return &c.Component, nil
}

// onMessageCreate listens for new messages and replies with
// "Ping!" or "Pong!" depending on the received message.
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID || m.Author.Bot {
        return
    }

    if m.Content == "ping" {
        _, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
        if nil != err {
            log.Warn().Msgf("Failed to deliver \"Pong!\" message: %v", err.Error())
        }
    }

    if m.Content == "pong" {
        _, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
        if nil != err {
            log.Warn().Msgf("Failed to deliver \"Ping!\" message: %v", err.Error())
        }
    }
}
