package tmuxt

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
)

type tmuxWrapperTestSuite struct {
}

func (c tmuxWrapperTestSuite) testTmuxWrapper(t *testing.T) {
	config := &Config{
		SessionName: gofakeit.LetterN(10),
		Windows: []*Window{
			{
				Name: gofakeit.LetterN(10),
				Grid: "vim vim vim term\nvim vim vim term\nplay play play play",
			},
			{
				Name: gofakeit.LetterN(10),
				Grid: "vim1 vim2 vim3\nvim1 vim2 vim3",
			},
			{
				Name: gofakeit.LetterN(10),
				Grid: "vim1",
			},
			{
				Name: gofakeit.LetterN(10),
				Grid: "vim1 vim2 vim3\nvim1 vim2 vim3\nvim1 vim2 vim3",
			},
		},
	}
	err := config.Validate()
	require.NoError(t, err)
	err = config.Parse()
	require.NoError(t, err)
	wrapper := NewTmuxWrapper(config, NewDimensions(274, 81))
	err = wrapper.Apply()
	require.NoError(t, err)
}

func TestTmuxWrapper(t *testing.T)  {
	suite := tmuxWrapperTestSuite{}

	t.Run("TmuxWrapper", suite.testTmuxWrapper)

}