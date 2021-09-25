package chaakoo

import (
	"fmt"
	"strings"
)

import (
	"errors"
)

// Config holds the entire config
type Config struct {
	SessionName string    `mapstructure:"name"`
	Windows     []*Window `mapstructure:"windows"`
	DryRun      bool
	ExitOnError bool
}

// Validate validates the config
// The error messages contain contextual information related to the validation issues
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if len(c.SessionName) == 0 {
		return errors.New("session name is required")
	}
	if len(c.Windows) == 0 {
		return fmt.Errorf("atleast 1 window is required for session - %s", c.SessionName)
	}
	for _, window := range c.Windows {
		if err := window.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Parse delegates to Window.Parse
func (c *Config) Parse() error {
	for _, window := range c.Windows {
		if err := window.Parse(); err != nil {
			return fmt.Errorf("unable to parse grid for window - %s: %w", window.Name, err)
		}
	}
	return nil
}

// Window represents one TMUX window from the config
type Window struct {
	Name      string `mapstructure:"name"`
	Grid      string `mapstructure:"grid"`
	FirstPane *Pane
	Commands  []*Command `mapstructure:"commands"`
}

// Validate validates a Window related config
func (w *Window) Validate() error {
	if w == nil {
		return errors.New("window is nil")
	}
	if len(w.Name) == 0 {
		return errors.New("window name is required")
	}
	if len(strings.TrimSpace(w.Grid)) == 0 {
		return fmt.Errorf("grid for window, %s, is empty", w.Name)
	}
	return nil
}

// Parse - parses the config
func (w *Window) Parse() error {
	if w == nil {
		return errors.New("window is nil")
	}
	grid, err := PrepareGrid(w.Grid)
	if err != nil {
		return err
	}
	pane, err := PrepareGraph(grid)
	if err != nil {
		return err
	}
	w.FirstPane = pane
	return nil
}

// Command represents a command fragment that will be executed in the pane whose name will be same as name in this
// struct.
// WorkingDirectory is the location in which all the commands will be executed.
// The working directory can be passed to tmux split-window command with -c flag but doing that will not create the
// pane if the working directory is wrong. So, in this implementation, passing the working directory is deferred until
// the pane has been created.
type Command struct {
	Name             string `mapstructure:"pane"`
	CommandText      string `mapstructure:"command"`
	WorkingDirectory string `mapstructure:"workdir"`
}
