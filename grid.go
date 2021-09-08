package tmuxt

import (
	"errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

var InvalidDimensionError = errors.New("invalid grid found! all rows must have same number of columns and vice versa")

func pregareGrid(gridKey string) ([][]string, error) {
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
	xStart uint
	xEnd   uint
	yStart uint
	yEnd   uint
    width  uint
    height uint
}


func prepareGraph(grid [][]string) {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
    //        currentGrid := grid[i][j]

		}
	}
}


func findBoundary(paneName string, startI, startJ int,  grid string[][]) {
    width := findWidth(paneName, startI, startJ, grid)
    height := findHeight(paneName, startI, startJ, grid)

    for i := startI; i < startI + height; i++ {
            for j := startJ; i < startJ + width; j++ {
                if grid[i][j] != paneName {
                    return 0, 0, errors.New(fmt.Sprintf("pane -> %s must be present at index %d, %d to make a rectangle", paneName, i, j))
                }
            }
    }

    return height, width, nil
}

func findWidth(paneName string, startI, startJ int, grid string[][]) (uint, error) {
    if (paneName != grid[i][j]) {
        return 0, errors.New(fmt.Sprintf("invalid pane, %s, for indexes %d, %d", paneName, startI, startJ))
    }

    var width uint := 0
    for col := startJ; col < len(grid[0]); col++ {
        if paneName == grid[startI][col] {
            width++
        }
    }

    return width, nil
}


func findHeight(paneName string, startI, startJ int, grid string[][]) (uint, error) {
    if (paneName != grid[i][j]) {
        return 0, errors.New(fmt.Sprintf("invalid pane, %s, for indexes %d, %d", paneName, startI, startJ))
    }

    var height uint := 0
    for row := startI; row< len(grid); row++ {
        if paneName == grid[row][startJ] {
            height++
        }
    }

    return height, nil
}
