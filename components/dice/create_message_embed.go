package dice

import "github.com/bwmarrin/discordgo"

// Create message embed with one field
func createMessageEmbed(title string, fieldTitle string, fieldContent string) discordgo.MessageEmbed {
	f := createMessageEmbedField(fieldTitle, fieldContent)
	fa := [1]*discordgo.MessageEmbedField{&f}
	e := discordgo.MessageEmbed{
		Title:  title,
		Fields: fa[:],
	}

	return e
}

// Create message embed field with a name and a value
func createMessageEmbedField(n string, v string) discordgo.MessageEmbedField {
	f := discordgo.MessageEmbedField{
		Name:  n,
		Value: v,
	}

	return f
}
