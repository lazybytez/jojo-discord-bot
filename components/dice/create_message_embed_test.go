package dice

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreateMessageEmbedTestSuite struct{ suite.Suite }

func TestCreateMessageEmbed(t *testing.T) {
	suite.Run(t, new(CreateMessageEmbedTestSuite))
}

func (suite *CreateMessageEmbedTestSuite) TestCreateAnswerEmbedMessage() {
	e := discordgo.MessageEmbed{
		Title: "You rolled 3 d6",
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

	suite.ObjectsAreEqual(e, r)
}

func (suite *CreateMessageEmbedTestSuite) ObjectsAreEqual(expected, actual interface{}) {
	if !assert.ObjectsAreEqual(expected, actual) {
		suite.T().Error("objects are not equal")
	}
}
