package dice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

var memberPermissions int64 = discordgo.PermissionSendMessages

var minValueDieSites = float64(2)
var maxValueDieSites = float64(1000)
var minValueDice = float64(1)
var maxValueDice = float64(100)
var diceCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:                     "dice",
		Description:              "throw one or more dice of your wished type.",
		DefaultMemberPermissions: &memberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "die-sites-number",
				Description: "The number of how many sites the die has, default is 6",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    &minValueDieSites,
				MaxValue:    maxValueDieSites,
				Required:    false,
			},
			{
				Name:        "number-dice",
				Description: "How many dice you want to throw, default is 1",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    &minValueDice,
				MaxValue:    maxValueDice,
				Required:    false,
			},
		},
	},
	Global:   true,
	Category: api.CategoryFun,
	Handler:  handleDice,
}

// handleDice handles the dice slash command
func handleDice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := getOptionsAsMap(i)
	n := getIntOption(options, "number-dice", 1)
	d := getIntOption(options, "die-sites-number", 6)

	r := rollDice(d, n)
	e := createAnswerEmbedMessage(n, d, r)
	eSlice := [1]*discordgo.MessageEmbed{&e}
	sendAnswer(s, i, eSlice[:])
}

// getOptionsAsMap create a map and insert the command options
func getOptionsAsMap(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

// getIntOption returns the option as int
func getIntOption(options map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, defaultValue int) int {
	o := defaultValue
	if opt, ok := options[name]; ok {
		o = int(opt.IntValue())
	}

	return o
}

// sendAnser sends the Answer
func sendAnswer(s *discordgo.Session, i *discordgo.InteractionCreate, e []*discordgo.MessageEmbed) {
	resp := &discordgo.InteractionResponseData{
		Embeds: e,
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
