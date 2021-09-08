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

func prepareGraph(grid [][]string) {

}
