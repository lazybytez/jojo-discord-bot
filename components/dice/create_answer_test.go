package dice

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CreateAnswerTestSuite struct{ suite.Suite }

func TestCreateAnswerEmbed(t *testing.T) {
	suite.Run(t, new(CreateAnswerTestSuite))
}

func (suite *CreateAnswerTestSuite) TestArrayIntToArrayString() {
	p := [5]int{4, 20, -13, 0, -6}
	e := [5]string{"4", "20", "-13", "0", "-6"}

	r := arrayIntToArrayString(p[:])

	suite.Equal(len(e), len(r))
	suite.ElementsMatch(e, r)
}

func (suite *CreateAnswerTestSuite) TestCreateAnswerResultContent() {
	a := [10]string{"What", "is", "'Courage'?", "Courage", "is", "owning", "your", "feat!", "-", "Zeppeli"}
	e := "What, is, 'Courage'?, Courage, is, owning, your, feat!, -, Zeppeli"

	r := createAnswerResultContent(a[:])

	suite.Equal(e, r)
}

func (suite *CreateAnswerTestSuite) TestCreateAnswerResultTitleOne() {
	e := "The Result is"

	r := createAnswerResultTitle(1)

	suite.Equal(e, r)
}

func (suite *CreateAnswerTestSuite) TestCreateAnswerResultTitleMultiple() {
	e := "The Results are"

	r := createAnswerResultTitle(2)

	suite.Equal(e, r)
}

func (suite *CreateAnswerTestSuite) TestCreateAnswerTitle() {
	e := "You rolled 3 d6"

	r := createAnswerTitle(3, 6)

	suite.Equal(e, r)
}
