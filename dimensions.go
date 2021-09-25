package chaakoo

import (
	"errors"
	"fmt"
	"golang.org/x/term"
)

var errNotInTerminalError = errors.New("to get the width and height, the program must run in the terminal")

// Dimension represents the dimension of the terminal in which binary is executed in.
type Dimension struct {
	Width  int
	Height int
}

// NewDimension constructs a dimension
func NewDimension(width int, height int) *Dimension {
	return &Dimension{Width: width, Height: height}
}

// TerminalDimension can be implemented by the providers of terminal dimensions
type TerminalDimension interface {
	Dimension() (*Dimension, error)
}

// DimensionUsingTerm uses golang/x/term to find the dimensions of the current shell.
// There can be other implementations using:
//	- tput cols and tput size
//	- stty size
type DimensionUsingTerm struct {
}

// Dimension fails if the binary is being executed in a terminal
func (d *DimensionUsingTerm) Dimension() (*Dimension, error) {
	if !term.IsTerminal(0) {
		return nil, errNotInTerminalError
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return nil, fmt.Errorf("cannot get the width and height: %w", err)
	}
	return NewDimension(width, height), nil
}
