package tmuxt

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

type GridTestCase struct {
	Error      bool
	GridActual [][]string
	Grid       string
}

func (g GridSuite) testPrepareGrid(t *testing.T) {
	var gridTestCases []GridTestCase
	if err := viper.UnmarshalKey("grids", &gridTestCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for i := range gridTestCases {
		result, err := PrepareGrid(gridTestCases[i].Grid)
		if gridTestCases[i].Error {
			require.Error(t, err)
			require.Nil(t, result)
		} else {
			require.NoError(t, err)
			require.ElementsMatch(t, gridTestCases[i].GridActual, result)
		}
	}
}
