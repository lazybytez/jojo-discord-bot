package dice

import (
	"net/http"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/test/discordgo_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct{ suite.Suite }

func TestCommand(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestHandleDice() {
	i := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:    "interaction.ID",
			Token: "interaction.Token",
			Data: discordgo.ApplicationCommandInteractionData{
				Options: createOptionArray(getDiceTestArray()),
			},
			Type: discordgo.InteractionApplicationCommand,
		},
	}

	s, r := discordgo_mock.MockSession()

	requestInteractionResponse := &discordgo.InteractionResponse{}
	r.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	handleDice(s, i)

	r.AssertExpectations(suite.T())

	checkResponseDataHasOneEmbedAndIsNotNil(suite, *requestInteractionResponse.Data)
	suite.ObjectsAreEqualValues(getDiceEmbedMessage(), *(requestInteractionResponse.Data.Embeds[0]))
}

func getDiceTestArray() map[string]float64 {
	t := make(map[string]float64)
	t["number-dice"] = 3
	t["die-sites-number"] = 1

	return t
}

func getDiceEmbedMessage() discordgo.MessageEmbed {
	m := discordgo.MessageEmbed{
		Title: "You rolled 3 d1",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "The Results are",
				Value:  "1, 1, 1",
				Inline: false,
			},
		},
	}

	return m
}

func (suite *CommandTestSuite) TestGetOptionsAsMap() {
	i := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Data: discordgo.ApplicationCommandInteractionData{
				Options: createOptionArray(getJoestarTestArray()),
			},
			Type: discordgo.InteractionApplicationCommand,
		},
	}
	expected := createOptionMap()

	actual := getOptionsAsMap(i)

	suite.Len(actual, len(expected))
	for i, obj := range expected {
		suite.ObjectsAreEqualValues(obj, actual[i])
	}
}

func (suite *CommandTestSuite) ObjectsAreEqualValues(expected, actual interface{}) {
	if !assert.ObjectsAreEqualValues(expected, actual) {
		suite.T().Error("objects are not equal")
	}
}

func (suite *CommandTestSuite) TestGetIntOptionGetValue() {
	optionMap := createOptionMap()

	intValue := getIntOption(optionMap, "Joseph Joestar", 3)

	suite.Equal(2, intValue)
}

func (suite *CommandTestSuite) TestGetIntOptionGetDefaultValue() {
	optionMap := createOptionMap()

	intValue := getIntOption(optionMap, "Dio Brando", 3)

	suite.Equal(3, intValue)
}

func (suite *CommandTestSuite) TestSendAnswer() {
	s, r := discordgo_mock.MockSession()
	i := discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:    "interaction.ID",
			Token: "interaction.Token",
		},
	}
	e := make([]*discordgo.MessageEmbed, 1)

	requestInteractionResponse := &discordgo.InteractionResponse{}
	r.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	sendAnswer(s, &i, e)

	r.AssertExpectations(suite.T())
	checkResponseDataHasOneEmbedAndIsNotNil(suite, *requestInteractionResponse.Data)
}

func checkResponseDataHasOneEmbedAndIsNotNil(s *CommandTestSuite, data discordgo.InteractionResponseData) {
	s.NotNil(data)
	s.NotNil(data.Embeds)
	s.Len(data.Embeds, 1)
}

func createOptionMap() map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, 3)

	for name, value := range getJoestarTestArray() {
		optionMap[name] = createDiscordOptionWithValue(name, value)
	}

	return optionMap
}

func createOptionArray(a map[string]float64) []*discordgo.ApplicationCommandInteractionDataOption {
	i := 0
	options := make([]*discordgo.ApplicationCommandInteractionDataOption, len(a))
	for name, value := range a {
		options[i] = createDiscordOptionWithValue(name, value)
		i++
	}

	return options
}

func getJoestarTestArray() map[string]float64 {
	t := make(map[string]float64)
	t["Jonathan Joestar"] = 1
	t["Joseph Joestar"] = 2
	t["Jotarou Joestar"] = 3

	return t
}

func createDiscordOptionWithValue(name string, value float64) *discordgo.ApplicationCommandInteractionDataOption {
	option := new(discordgo.ApplicationCommandInteractionDataOption)
	option.Name = name
	option.Value = value
	option.Type = discordgo.ApplicationCommandOptionInteger

	return option
}
