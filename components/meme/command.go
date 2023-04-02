/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2023 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package meme

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// memeCommand is the command for posting memes.
// It allows to post memes from different sources.
var memeCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "meme",
		Description: "Post some memes into the current channel!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "subreddit",
				Description: "The subreddit to post from",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices:     memeSubreddits,
			},
		},
	},
	Category: api.CategoryFun,
	Handler:  handleMeme,
}

// handleMeme handles the meme command.
// It checks whether the root command or a subcommand was used and then calls the corresponding handler.
// If no subcommand was used, it will post a meme from a random subreddit.
func handleMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	randomMeme, err := findRandomMeme("", 0)

	if err != nil {
		responseTemplate := slash_commands.GenerateEphemeralInteractionResponseTemplate("", "")

		responseTemplate.Embeds = []*discordgo.MessageEmbed{
			{
				Title:       "Error",
				Description: "An error occurred while trying to post a meme.",
			},
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: responseTemplate,
		})

		user := i.User
		if nil == user {
			user = i.Member.User
		}

		C.Logger().Err(err, fmt.Sprintf(
			"An error occurred while trying to post a meme on guild \"%s\" as requested by user \"%s\" / \"%s:%s\".",
			i.GuildID, user.ID, user.Username, user.Discriminator))

		return
	}

	responseTemplate := slash_commands.GenerateEphemeralInteractionResponseTemplate("", "")
	responseTemplate.Embeds = []*discordgo.MessageEmbed{
		{
			Title:       randomMeme.Title,
			Description: fmt.Sprintf("Posted on %s", randomMeme.Subreddit),
			URL:         randomMeme.URL,
			Image: &discordgo.MessageEmbedImage{
				URL: randomMeme.ImageURL,
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseTemplate,
	})
}
