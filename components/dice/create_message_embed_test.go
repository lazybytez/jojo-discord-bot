package dice

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestCreateAnswerEmbedMessage(t *testing.T) {
	e := discordgo.MessageEmbed{
		URL:         "",
		Type:        "",
		Title:       "You rolled 3 d6",
		Description: "",
		Timestamp:   "",
		Color:       0,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "The Results are",
				Value:  "3, 5, 2",
				Inline: false,
			},
		},
	}
	a := [3]int{3, 5, 2}

	r := createAnswerEmbedMessage(3, 6, a[:])

	checkIfTwoMessageEmbedAreTheSame(t, e, r)
}

func checkIfTwoMessageEmbedAreTheSame(t *testing.T, e discordgo.MessageEmbed, g discordgo.MessageEmbed) {
	checkIfTwoStringsAreTheSame(t, e.URL, g.URL, "urls")
	checkIfTwoStringsAreTheSame(t, string(e.Type), string(g.Type), "Types")
	checkIfTwoStringsAreTheSame(t, e.Title, g.Title, "titles")
	checkIfTwoStringsAreTheSame(t, e.Description, g.Description, "descriptions")
	checkIfTwoStringsAreTheSame(t, e.Timestamp, g.Timestamp, "timestamps")
	checkIfTwoIntAreTheSame(t, e.Color, g.Color, "colors")
	checkIfTwoIntAreTheSame(t, len(e.Fields), len(g.Fields), "number fields")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Footer, g.Footer, "footers")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Image, g.Image, "images")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Thumbnail, g.Thumbnail, "thumbnails")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Video, g.Video, "videos")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Provider, g.Provider, "providers")
	checkIfTwoObjectsAreTheSameInDepth(t, e.Author, g.Author, "authors")
	checkIfTwoMessageEmbedFieldsSlicesAreTheSame(t, e.Fields, g.Fields)
	checkIfTwoObjectsAreTheSameInDepth(t, e, g, "message-embed")
}

func checkIfTwoMessageEmbedFieldsSlicesAreTheSame(t *testing.T, e []*discordgo.MessageEmbedField, g []*discordgo.MessageEmbedField) {
	for k, elem := range e {
		gElem := g[k]
		checkIfTwoMessageEmbedFieldsAreTheSame(t, *elem, *gElem)
	}
}

func checkIfTwoMessageEmbedFieldsAreTheSame(t *testing.T, e discordgo.MessageEmbedField, g discordgo.MessageEmbedField) {
	checkIfTwoStringsAreTheSame(t, e.Name, g.Name, "names")
	checkIfTwoStringsAreTheSame(t, e.Value, g.Value, "values")
	checkIfTwoBoolAreTheSame(t, e.Inline, g.Inline, "inlines")
}
