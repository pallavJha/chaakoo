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
    parent *pane
}

func prepareGraph(grid [][]string) {
    panes, err := preparePanes(grid)
    if err != nil {
        // handle error
    }

    



}


func getLeftPaneName(pane *pane, grid [][]string, panes map[string]*pane) string {
    i, j := pane.yStart, pane.xEnd + 1
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
    i, j := pane.yEnd + 1, pane.xStart
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
