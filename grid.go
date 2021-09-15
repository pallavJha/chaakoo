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

func prepareGraph(grid [][]string) (*Pane, error) {
	panes, err := preparePanes(grid)
	if err != nil {
		log.Debug().Interface("grid", grid).Err(err).Msg("cannot convert grid to panes")
		return nil, err
	}

	var nameToPanes = make(map[string]*Pane)
	var visited = make(map[string]bool)
	for _, pane := range panes {
		if prevPane, ok := nameToPanes[pane.Name]; !ok {
			nameToPanes[pane.Name] = pane
			visited[pane.Name] = false
		} else {
			log.Debug().
				Str("grid",
					fmt.Sprintf("(%d, %d), (%d, %d)", prevPane.XStart, prevPane.YStart, prevPane.XEnd, prevPane.YEnd),
				).
				Str("grid",
					fmt.Sprintf("(%d, %d), (%d, %d)", pane.XStart, pane.YStart, pane.XEnd, pane.YEnd),
				).
				Msg("pane, %s, appears multiple times")
		}
	}

	numberOfPanes := len(panes)
	for numberOfPanes > 0 {
		dfs(panes[0], grid, nameToPanes, visited)
		if !continueDFS(panes[0], nameToPanes) {
			break
		}
		numberOfPanes--
	}

	if continueDFS(panes[0], nameToPanes) {
		log.Debug().Interface("grid", grid).Msg("cannot create a pane arrangement for the provided grid after dfs traversal")
		return nil, errors.New("cannot create a pane arrangement for the provided grid after dfs traversal")
	}

	return panes[0], err
}

func dfs(currentPane *Pane, grid [][]string, panes map[string]*Pane, visited map[string]bool) *Pane {
	leftPaneName := getLeftPaneName(currentPane, grid, panes)
	var leftPane *Pane
	if len(leftPaneName) > 0 {
		leftPane = panes[leftPaneName]
		leftPane = dfs(leftPane, grid, panes, visited)
	}

	bottomPaneName := getBottomPaneName(currentPane, grid, panes)
	var bottomPane *Pane
	if len(bottomPaneName) > 0 {
		bottomPane = panes[bottomPaneName]
		bottomPane = dfs(bottomPane, grid, panes, visited)
	}

	if leftPane != nil && bottomPane == nil {
		if leftPane.Height() == currentPane.Height() {
			// directly add the left pane as they have the same height
			currentPane.AddLeftPane(leftPane)
			currentPane.XEnd = leftPane.XEnd
			leftPane.Visited = true
			visited[leftPaneName] = true
		}
	} else if leftPane == nil && bottomPane != nil {
		if bottomPane.Width() == currentPane.Width() {
			// directly add the bottom pane as they have the same width
			currentPane.AddBottomPane(bottomPane)
			currentPane.YEnd = bottomPane.YEnd
			bottomPane.Visited = true
			visited[bottomPaneName] = true
		}
	} else if leftPane != nil && bottomPane != nil {
		if leftPane.Height() == currentPane.Height() {
			// directly add the left pane as they have the same height
			currentPane.AddLeftPane(leftPane)
			currentPane.XEnd = leftPane.XEnd
			leftPane.Visited = true
		} else if leftPane.Height() > currentPane.Height() {
			// here the left pane has more height than the currentPane,
			// so we join the bottom pane to the current pane to increase
			// the height and
			currentPane.AddBottomPane(bottomPane)
			currentPane.YEnd = bottomPane.YEnd
			bottomPane.Visited = true
			// then join the left pane
			currentPane.AddLeftPane(leftPane)
			currentPane.XEnd = leftPane.XEnd
			leftPane.Visited = true
		}
	}
	return currentPane
}

func continueDFS(firstPane *Pane, panes map[string]*Pane) bool {
	for _, pane := range panes {
		if pane.Name != firstPane.Name {
			if !pane.Visited {
				return true
			}
		}
	}

	return false
}

func getLeftPaneName(pane *Pane, grid [][]string, panes map[string]*Pane) string {
	i, j := pane.YStart, pane.XEnd+1
	if j >= len(grid[0]) {
		return ""
	}
	leftPaneName := grid[i][j]
	leftPane := panes[leftPaneName]
	if leftPane.YStart == pane.YStart && !leftPane.Visited {
		return leftPaneName
	}
	return ""
}

func getBottomPaneName(pane *Pane, grid [][]string, panes map[string]*Pane) string {
	i, j := pane.YEnd+1, pane.XStart
	if i >= len(grid) {
		return ""
	}
	bottomPaneName := grid[i][j]
	bottomPane := panes[bottomPaneName]
	if bottomPane.XStart == pane.XStart && !bottomPane.Visited {
		return bottomPaneName
	}
	return ""
}

func preparePanes(grid [][]string) ([]*Pane, error) {
	var visited = make([][]bool, len(grid))
	for i := 0; i < len(grid); i++ {
		visited[i] = make([]bool, len(grid[0]))
	}

	var panes []*Pane

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if !visited[i][j] {
				height, width, err := findBoundary(grid[i][j], i, j, grid, visited)
				if err != nil {
					return nil, err
				}
				panes = append(panes, &Pane{
					Name:   grid[i][j],
					XStart: j,
					YStart: i,
					XEnd:   j + width - 1,
					YEnd:   i + height - 1,
				})
			}
		}
	}

	return panes, nil
}

func findBoundary(paneName string, startI, startJ int, grid [][]string, visited [][]bool) (int, int, error) {
	width := findWidth(paneName, startI, startJ, grid)
	height := findHeight(paneName, startI, startJ, grid)

	for i := startI; i < startI+height; i++ {
		for j := startJ; j < startJ+width; j++ {
			if grid[i][j] != paneName {
				return 0, 0, errors.New(fmt.Sprintf("pane -> %s must be present at index %d, %d to make a rectangle", paneName, i, j))
			} else {
				visited[i][j] = true
			}
		}
	}

	return height, width, nil
}

func findWidth(paneName string, startI, startJ int, grid [][]string) int {
	var width = 0
	for col := startJ; col < len(grid[0]); col++ {
		if paneName == grid[startI][col] {
			width++
		}
	}

	return width
}

func findHeight(paneName string, startI, startJ int, grid [][]string) int {
	var height = 0
	for row := startI; row < len(grid); row++ {
		if paneName == grid[row][startJ] {
			height++
		}
	}

	return height
}
