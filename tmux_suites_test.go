package tmuxt

import "testing"

type GridSuite struct {
}

func TestGrid(t *testing.T) {
	suite := GridSuite{}
	t.Run("TestPrepareGrid", suite.testPrepareGrid)
}