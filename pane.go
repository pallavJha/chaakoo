package chaakoo

// Pane represents a TMUX pane in a 2D grid
type Pane struct {
	Name             string  // Name of the pane
	XStart           int     // First index in the horizontal direction
	XEnd             int     // Last index in the horizontal direction
	YStart           int     // First index in the vertical direction
	YEnd             int     // Last index in the vertical direction
	Visited          bool    // If this node was visited while traversal
	Left             []*Pane // Collection of the panes to the left of the current pane
	Bottom           []*Pane // Collection of the panes to the bottom of the current pane
	priorLeftIndex   int     // used while dfs
	priorBottomIndex int     // used while dfs
}

// Height returns the height of the pane
func (p *Pane) Height() int {
	return p.YEnd - p.YStart + 1
}

// Width returns the width of the pane
func (p *Pane) Width() int {
	return p.XEnd - p.XStart + 1
}

// AddLeftPane appends left pane to the current pane
func (p *Pane) AddLeftPane(leftPane *Pane) {
	p.Left = append(p.Left, leftPane)
}

// AddBottomPane appends a bottom pane to the current pane
func (p *Pane) AddBottomPane(bottomPane *Pane) {
	p.Bottom = append(p.Bottom, bottomPane)
}

func (p *Pane) reset() {
	p.priorLeftIndex = len(p.Left) - 1
	p.priorBottomIndex = len(p.Bottom) - 1
}

// AsGrid returns the 2D string array representation of the Pane
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
		} else if leftPane.Height() > bottomPane.Height() {
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
