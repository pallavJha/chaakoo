package chaakoo

import (
	"errors"
	"fmt"
	"golang.org/x/term"
)

var NotInTerminalError = errors.New("to get the width and height, the program must run in the terminal")

type Dimensions struct {
	Width  int
	Height int
}

func NewDimensions(width int, height int) *Dimensions {
	return &Dimensions{Width: width, Height: height}
}

type TerminalDimensions interface {
	Dimensions() (*Dimensions, error)
}

type DimensionsUsingTerm struct {
}

func (d *DimensionsUsingTerm) Dimensions() (*Dimensions, error) {
	if !term.IsTerminal(0) {
		return nil, NotInTerminalError
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return nil, fmt.Errorf("cannot get the width and height: %w", err)
	}
	return NewDimensions(width, height), nil
}
