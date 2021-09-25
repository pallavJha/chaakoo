package chaakoo

import (
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

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
	ci := os.Getenv("CI")
	log.Info().Str("CI", ci).Msg("ci environment check")
	if len(ci) > 0 && ci == "true" {
		readTestConfig("tmux_wrapper_apply_test_cases_ci")
	} else {
		readTestConfig("tmux_wrapper_apply_test_cases")
	}
	t.Run("TmuxWrapper", suite.testTmuxWrapperApply)

}
