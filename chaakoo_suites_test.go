package chaakoo

import "testing"

type GridSuite struct {
}

func TestPrepareGrid(t *testing.T) {
	suite := GridSuite{}
	readTestConfig("prepare_grid_testcases")
	t.Run("TestPrepareGrid", suite.testPrepareGrid)
}

func TestPreparePane(t *testing.T) {
	suite := GridSuite{}
	readTestConfig("prepare_pane_testcases")
	t.Run("TestPreparePane", suite.testPreparePanes)
}

func TestPrepareGraph(t *testing.T) {
	suite := GridSuite{}
	readTestConfig("prepare_graph_testcases")
	t.Run("TestPrepareGraph", suite.testPrepareGraph)
}

type TmuxWrapperTestSuite struct {

}

func TestTmuxWrapper_Apply(t *testing.T)  {
	suite := TmuxWrapperTestSuite{}
	readTestConfig("tmux_wrapper_apply_test_cases")
	t.Run("TmuxWrapper", suite.testTmuxWrapperApply)

}
