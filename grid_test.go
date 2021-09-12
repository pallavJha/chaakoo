package tmuxt

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

func (g GridSuite) testPrepareGraph1(t *testing.T) {
	grid, err := PrepareGrid(`a a b
						              c d e`)
	require.NoError(t, err)
	a, err := prepareGraph(grid)
	require.NoError(t, err)

	require.NotNil(t, a)
	matchPaneAttributes(t, a, "a", 0, 2, 0, 1)

	matchPaneAttributes(t, a.Left, "b", 2, 2, 0, 1)
	require.Nil(t, a.Left.Left)
	matchPaneAttributes(t, a.Left.Bottom, "e", 2, 2, 1, 1)

	matchPaneAttributes(t, a.Bottom, "c", 0, 1, 1, 1)
	require.Nil(t, a.Bottom.Bottom)
	matchPaneAttributes(t, a.Bottom.Left, "d", 1, 1, 1, 1)
}

func (g GridSuite) testPrepareGraph2(t *testing.T) {
	grid, err := PrepareGrid(`a b b
						              c d e
						              e f e`)
	require.NoError(t, err)
	a, err := prepareGraph(grid)
	require.NoError(t, err)

	require.NotNil(t, a)
	matchPaneAttributes(t, a, "a", 0, 2, 0, 1)

	matchPaneAttributes(t, a.Left, "b", 2, 2, 0, 1)
	require.Nil(t, a.Left.Left)
	matchPaneAttributes(t, a.Left.Bottom, "e", 2, 2, 1, 1)

	matchPaneAttributes(t, a.Bottom, "c", 0, 1, 1, 1)
	require.Nil(t, a.Bottom.Bottom)
	matchPaneAttributes(t, a.Bottom.Left, "d", 1, 1, 1, 1)
}

func matchPaneAttributes(t *testing.T, pane *Pane, name string, xStart, xEnd, yStart, yEnd int) {
	require.Equal(t, pane.Name, name)
	require.Equal(t, pane.XStart, xStart)
	require.Equal(t, pane.XEnd, xEnd)
	require.Equal(t, pane.YEnd, yEnd)
	require.Equal(t, pane.YEnd, yEnd)
}
