package chaakoo

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

type GridTestCase struct {
	ID         int
	Error      bool
	GridActual [][]string
	Grid       string
	Panes      []*Pane
	PaneError  string
}

func (g GridSuite) testPrepareGrid(t *testing.T) {
	var gridTestCases []GridTestCase
	if err := viper.UnmarshalKey("grids", &gridTestCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for i := range gridTestCases {
		t.Log("Test case", gridTestCases[i].ID)
		result, err := PrepareGrid(gridTestCases[i].Grid)
		if gridTestCases[i].Error {
			require.Error(t, err)
			require.Nil(t, result)
		} else {
			require.NoError(t, err)
			require.Equal(t, gridTestCases[i].GridActual, result)
		}
	}
}

func (g GridSuite) testPreparePanes(t *testing.T) {
	var gridTestCases []GridTestCase
	if err := viper.UnmarshalKey("grids", &gridTestCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for _, testcase := range gridTestCases {
		if !testcase.Error {
			t.Log("Test case", testcase.ID)
			grid, err := PrepareGrid(testcase.Grid)
			result, err := preparePanes(grid)
			if len(testcase.PaneError) == 0 {
				require.NoError(t, err)

				var expectedPanes = make(map[string]*Pane)
				for _, pane := range testcase.Panes {
					expectedPanes[pane.Name] = pane
				}
				var actualPanes = make(map[string]*Pane)
				for _, pane := range result {
					actualPanes[pane.Name] = pane
				}

				require.Equal(t, len(expectedPanes), len(actualPanes))
				for paneName := range expectedPanes {
					require.Equal(t, expectedPanes[paneName], actualPanes[paneName])
				}
			} else {
				require.Error(t, err)
				require.EqualError(t, err, testcase.PaneError)
			}
		}
	}
}

func (g GridSuite) testPrepareGraph(t *testing.T) {
	var gridTestCases []GridTestCase
	if err := viper.UnmarshalKey("grids", &gridTestCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for _, testcase := range gridTestCases {
		t.Log("Test case", testcase.ID)
		if testcase.ID == 7 {
			t.Log("For Debug")
		}
		grid, err := PrepareGrid(testcase.Grid)
		require.NoError(t, err)
		topPane, err := PrepareGraph(grid)
		if !testcase.Error {
			require.NoError(t, err)
			require.Equal(t, topPane.AsGrid(), testcase.GridActual)
		} else {
			require.Error(t, err)
			require.EqualError(t, err, testcase.PaneError)
		}
	}
}
