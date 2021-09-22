package chaakoo

type Pane struct {
	Name             string
	XStart           int
	XEnd             int
	YStart           int
	YEnd             int
	Visited          bool
	Left             []*Pane
	Bottom           []*Pane
	priorLeftIndex   int
	priorBottomIndex int
}

func (p *Pane) Height() int {
	return p.YEnd - p.YStart + 1
}

func (p *Pane) Width() int {
	return p.XEnd - p.XStart + 1
}

func (p *Pane) AddLeftPane(leftPane *Pane) {
	p.Left = append(p.Left, leftPane)
}

func (p *Pane) AddBottomPane(bottomPane *Pane) {
	p.Bottom = append(p.Bottom, bottomPane)
}

func (p *Pane) reset() {
	p.priorLeftIndex = len(p.Left) - 1
	p.priorBottomIndex = len(p.Bottom) - 1
}

func (p *Pane) AsGrid() [][]string {
	p.reset()
	var grid = make([][]string, p.Height())
	for i := range grid {
		grid[i] = make([]string, p.Width())
		for j := range grid[i] {
			grid[i][j] = p.Name
		}
	}
	for {
		var leftPane, bottomPane *Pane
		if p.priorLeftIndex > -1 {
			leftPane = p.Left[p.priorLeftIndex]
		}
		if p.priorBottomIndex > -1 {
			bottomPane = p.Bottom[p.priorBottomIndex]
		}
		if leftPane == nil && bottomPane == nil {
			return grid
		} else if leftPane != nil && bottomPane == nil {
			p.priorLeftIndex--
			leftGrid := leftPane.AsGrid()
			fill(p, leftPane, grid, leftGrid)
		} else if leftPane == nil && bottomPane != nil {
			p.priorBottomIndex--
			bottomGrid := bottomPane.AsGrid()
			fill(p, bottomPane, grid, bottomGrid)
		} else if leftPane.Height() > bottomPane.Width() {
			p.priorLeftIndex--
			leftGrid := leftPane.AsGrid()
			fill(p, leftPane, grid, leftGrid)
		} else {
			p.priorBottomIndex--
			bottomGrid := bottomPane.AsGrid()
			fill(p, bottomPane, grid, bottomGrid)
		}
	}
}

func fill(parent *Pane, child *Pane, parentGrid [][]string, childGrid [][]string) {
	for i, I := child.YStart-parent.YStart, 0; i <= child.YEnd-parent.YStart; i, I = i+1, I+1 {
		for j, J := child.XStart-parent.XStart, 0; j <= child.XEnd-parent.XStart; j, J = j+1, J+1 {
			parentGrid[i][j] = childGrid[I][J]
		}
	}
}
