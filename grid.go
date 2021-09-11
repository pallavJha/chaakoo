package tmuxt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var InvalidDimensionError = errors.New("invalid grid found! all rows must have same number of columns and vice versa")

func PrepareGrid(gridKey string) ([][]string, error) {
	gridKey = strings.TrimSpace(gridKey)
	var re = regexp.MustCompile(`\s+`)
	gridLines := strings.Split(gridKey, "\n")

	var grid [][]string
	for _, gridLine := range gridLines {
		gridLine = strings.TrimSpace(gridLine)
		gridLine = re.ReplaceAllString(gridLine, " ")
		gridCells := strings.Split(gridLine, " ")
		if len(gridCells) > 0 {
			var cells []string
			for _, cell := range gridCells {
				cells = append(cells, cell)
			}
			grid = append(grid, cells)
		}
	}

	if !checkForEqualWidth(grid) {
		log.Debug().Interface("grid", grid).Msg(InvalidDimensionError.Error())
		return nil, InvalidDimensionError
	}
	return grid, nil
}

func checkForEqualWidth(grid [][]string) bool {
	numberOfCellsInARow := len(grid[0])
	for _, row := range grid {
		if len(row) != numberOfCellsInARow {
			return false
		}
	}
	return true
}

type pane struct {
	name    string
	xStart  int
	xEnd    int
	yStart  int
	yEnd    int
	visited bool
	left    *pane
	botton  *pane
	parent  *pane
}

func (p *pane) Height() int {
	return p.yEnd - p.yStart + 1
}

func (p *pane) Width() int {
	return p.xEnd - p.xStart + 1
}

func prepareGraph(grid [][]string) {
	panes, err := preparePanes(grid)
	if err != nil {
		// handle error
	}

	var nameToPanes = make(map[string]*pane)
	var visited = make(map[string]bool)
	for _, pane := range panes {
		nameToPanes[pane.name] = pane
		visited[pane.name] = false
	}

  numberOfPanes := len(panes)
  for numberOfPanes > 0 {
    dfs(panes[0], grid, nameToPanes, visited)
    numberOfPanes--
  }

}

func dfs(currentPane *pane, grid [][]string, panes map[string]*pane, visited map[string]bool) *pane {
	leftPaneName := getLeftPaneName(currentPane, grid, panes)
	var leftPane *pane
	if len(leftPaneName) > 0 {
		leftPane = panes[leftPaneName]
		leftPane = dfs(leftPane, grid, panes, visited)
	}

	bottomPaneName := getBottomPaneName(currentPane, grid, panes)
	var bottomPane *pane
	if len(bottomPaneName) > 0 {
		bottomPane = panes[bottomPaneName]
		bottomPane = dfs(bottomPane, grid, panes, visited)
	}

	if leftPane != nil && bottomPane == nil {
		if leftPane.Height() == currentPane.Height() {
			// directly add the left pane as they have the same height
			currentPane.left = leftPane
			currentPane.xEnd = leftPane.xEnd
			leftPane.visited = true
			visited[leftPaneName] = true
		}
	} else if leftPane == nil && bottomPane != nil {
		if bottomPane.Width() == currentPane.Width() {
			// directly add the right pane as they have the same width
			currentPane.botton = bottomPane
			currentPane.yEnd = bottomPane.yEnd
			bottomPane.visited = true
			visited[bottomPaneName] = true
		}
	} else if leftPane != nil && bottomPane != nil {
		if leftPane.Height() == currentPane.Height() {
			// directly add the left pane as they have the same height
			currentPane.left = leftPane
			currentPane.xEnd = leftPane.xEnd
			leftPane.visited = true
		} else if leftPane.Height() > currentPane.Height() {
			// here the left pane has more height than the currentPane
			// so we join the bottom pane to the current pane to increase
			// the height and
			currentPane.botton = bottomPane
			currentPane.yEnd = bottomPane.yEnd
			bottomPane.visited = true
			// then join the left pane
			currentPane.left = leftPane
			currentPane.xEnd = leftPane.xEnd
			leftPane.visited = true
		}
	}
	return currentPane
}

func getLeftPaneName(pane *pane, grid [][]string, panes map[string]*pane) string {
	i, j := pane.yStart, pane.xEnd+1
	if j >= len(grid[0]) {
		return ""
	}
	leftPaneName := grid[i][j]
	leftPane := panes[leftPaneName]
	if leftPane.yStart == pane.yStart {
		return leftPaneName
	}
	return ""
}

func getBottomPaneName(pane *pane, grid [][]string, panes map[string]*pane) string {
	i, j := pane.yEnd+1, pane.xStart
	if i >= len(grid) {
		return ""
	}
	bottomPaneName := grid[i][j]
	bottomPane := panes[bottomPaneName]
	if bottomPane.xStart == pane.xStart {
		return bottomPaneName
	}
	return ""
}

func preparePanes(grid [][]string) ([]*pane, error) {
	var visited [][]bool = make([][]bool, len(grid))
	for i := 0; i < len(grid); i++ {
		visited[i] = make([]bool, len(grid[0]))
	}

	var panes []*pane

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if !visited[i][j] {
				height, width, err := findBoundary(grid[i][j], i, j, grid, visited)
				if err != nil {
					return nil, err
				}
				panes = append(panes, &pane{
					name:   grid[i][j],
					xStart: j,
					yStart: i,
					xEnd:   j + height - 1,
					yEnd:   i + width - 1,
				})
			}
		}
	}

	return panes, nil
}

func findBoundary(paneName string, startI, startJ int, grid [][]string, visited [][]bool) (int, int, error) {
	width, err := findWidth(paneName, startI, startJ, grid)
	if err != nil {
		return 0, 0, err
	}
	height, err := findHeight(paneName, startI, startJ, grid)
	if err != nil {
		return 0, 0, nil
	}

	for i := startI; i < startI+height; i++ {
		for j := startJ; i < startJ+width; j++ {
			if grid[i][j] != paneName {
				return 0, 0, errors.New(fmt.Sprintf("pane -> %s must be present at index %d, %d to make a rectangle", paneName, i, j))
			} else {
				visited[i][j] = true
			}
		}
	}

	return height, width, nil
}

func findWidth(paneName string, startI, startJ int, grid [][]string) (int, error) {
	if paneName != grid[startI][startJ] {
		return 0, errors.New(fmt.Sprintf("invalid pane, %s, for indexes %d, %d", paneName, startI, startJ))
	}

	var width int = 0
	for col := startJ; col < len(grid[0]); col++ {
		if paneName == grid[startI][col] {
			width++
		}
	}

	return width, nil
}

func findHeight(paneName string, startI, startJ int, grid [][]string) (int, error) {
	if paneName != grid[startI][startJ] {
		return 0, errors.New(fmt.Sprintf("invalid pane, %s, for indexes %d, %d", paneName, startI, startJ))
	}

	var height int = 0
	for row := startI; row < len(grid); row++ {
		if paneName == grid[row][startJ] {
			height++
		}
	}

	return height, nil
}
