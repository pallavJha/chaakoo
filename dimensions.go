package chaakoo

import (
	"errors"
	"fmt"
	"golang.org/x/term"
)

var NotInTerminalError = errors.New("to get the width and height, the program must run in the terminal")

type Dimension struct {
	Width  int
	Height int
}

func NewDimension(width int, height int) *Dimension {
	return &Dimension{Width: width, Height: height}
}

type TerminalDimension interface {
	Dimension() (*Dimension, error)
}

type DimensionUsingTerm struct {
}

func (d *DimensionUsingTerm) Dimension() (*Dimension, error) {
	if !term.IsTerminal(0) {
		return nil, NotInTerminalError
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return nil, fmt.Errorf("cannot get the width and height: %w", err)
	}
	return NewDimension(width, height), nil
}
