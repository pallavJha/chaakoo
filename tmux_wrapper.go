package tmuxt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strconv"
	"strings"
)

var CommandName = "tmux"

type TmuxError struct {
	stdout, stderr string
	err            error
}

func NewTmuxError(stdout, stderr string, err error) *TmuxError {
	return &TmuxError{
		stdout: stdout,
		stderr: stderr,
		err:    err,
	}
}

func (t *TmuxError) Error() string {
	return fmt.Sprintf("err: %v, stdout: %s, stderr: %s", t.err, t.stdout, t.stderr)
}

type TmuxCmdResponse struct {
	SessionID string
	WindowID  string
	PaneID    string
}

type TmuxWrapper struct {
	config     *Config
	dimensions *Dimensions
	executor   ICommandExecutor
}

func NewTmuxWrapper(config *Config, dimensions *Dimensions) *TmuxWrapper {
	wrapper := &TmuxWrapper{
		config:     config,
		dimensions: dimensions,
	}
	if config.DryRun {
		wrapper.executor = NewNOOPExecutor()
	} else {
		wrapper.executor = NewCommandExecutor()
	}
	return wrapper
}

func (t *TmuxWrapper) Apply() error {
	res, err := t.newSession(t.config.SessionName, t.config.Windows[0].Name, t.dimensions)
	if err != nil {
		t.killSession(t.config.SessionName)
		return err
	}
	var paneNames = make(map[string]string)
	paneNames[t.config.Windows[0].FirstPane.Name] = res.PaneID
	if err = t.walkPane(t.config.Windows[0].FirstPane, paneNames); err != nil {
		t.killSession(t.config.SessionName)
		return err
	}
	for i := 1; i < len(t.config.Windows); i++ {
		res, err = t.newWindow(t.config.SessionName, t.config.Windows[i].Name)
		if err != nil {
			t.killSession(t.config.SessionName)
			return err
		}
		paneNames = make(map[string]string)
		paneNames[t.config.Windows[i].FirstPane.Name] = res.PaneID
		if err = t.walkPane(t.config.Windows[i].FirstPane, paneNames); err != nil {
			t.killSession(t.config.SessionName)
			return err
		}
	}
	return nil
}

func (t *TmuxWrapper) walkPane(currentPane *Pane, paneNames map[string]string) error {
	currentPane.reset()
	for {
		var leftPane, bottomPane *Pane
		if currentPane.priorLeftIndex > -1 {
			leftPane = currentPane.Left[currentPane.priorLeftIndex]
		}
		if currentPane.priorBottomIndex > -1 {
			bottomPane = currentPane.Bottom[currentPane.priorBottomIndex]
		}
		if leftPane == nil && bottomPane == nil {
			return nil
		} else if leftPane != nil && bottomPane == nil {
			currentPane.priorLeftIndex--
			sizeInPercentage := float64(leftPane.Width()*100) / float64(currentPane.Width())
			res, err := t.newPane(paneNames[currentPane.Name], int(sizeInPercentage), true)
			if err != nil {
				return err
			}
			paneNames[leftPane.Name] = res.PaneID
			err = t.walkPane(leftPane, paneNames)
			if err != nil {
				return err
			}
		} else if leftPane == nil && bottomPane != nil {
			currentPane.priorBottomIndex--
			sizeInPercentage := float64(bottomPane.Height()*100) / float64(currentPane.Height())
			res, err := t.newPane(paneNames[currentPane.Name], int(sizeInPercentage), false)
			if err != nil {
				return err
			}
			paneNames[bottomPane.Name] = res.PaneID
			err = t.walkPane(bottomPane, paneNames)
			if err != nil {
				return err
			}
		} else if leftPane.Height() > bottomPane.Width() {
			currentPane.priorLeftIndex--
			sizeInPercentage := float64(leftPane.Width()*100) / float64(currentPane.Width())
			res, err := t.newPane(paneNames[currentPane.Name], int(sizeInPercentage), true)
			if err != nil {
				return err
			}
			paneNames[leftPane.Name] = res.PaneID
			err = t.walkPane(leftPane, paneNames)
			if err != nil {
				return err
			}
		} else {
			currentPane.priorBottomIndex--
			sizeInPercentage := float64(bottomPane.Height()*100) / float64(currentPane.Height())
			res, err := t.newPane(paneNames[currentPane.Name], int(sizeInPercentage), false)
			if err != nil {
				return err
			}
			paneNames[bottomPane.Name] = res.PaneID
			err = t.walkPane(bottomPane, paneNames)
			if err != nil {
				return err
			}
		}
	}
}

func (t *TmuxWrapper) newSession(sessionID, windowName string, dimensions *Dimensions) (*TmuxCmdResponse, error) {
	// tmux new-session -d -s session2 -n vim -x 136 -y 80
	var args = []string{
		"new-session",
		"-d",
		"-s",
		sessionID,
		"-n",
		windowName,
		"-x",
		strconv.Itoa(dimensions.Width),
		"-y",
		strconv.Itoa(dimensions.Height),
		"-P",
		"-F",
		"\"#{window_id}--#{pane_id}\"",
	}
	stdout, stderr, err := t.executor.Execute(CommandName, args...)
	if err != nil {
		return nil, NewTmuxError(stdout, stderr, err)
	}
	output := stdout
	strings.TrimSpace(output)
	splitOutput := strings.Split(output, "--")
	if len(splitOutput) != 2 {
		log.Debug().Interface("output", splitOutput).Msg("invalid output from list-panes sub command")
		return nil, NewTmuxError(stdout, "", errors.New("cannot parse the windowID and pane ID from the list-panes output"))
	}
	return &TmuxCmdResponse{
		SessionID: "",
		WindowID:  splitOutput[0][1:],
		PaneID:    splitOutput[1][1:],
	}, nil
}

func (t *TmuxWrapper) newWindow(sessionID, windowName string) (*TmuxCmdResponse, error) {
	// tmux new-window -t session3 -n vim2  -P -F "#{window_id}--#{pane_id}"
	// @9--%19

	var args = []string{
		"new-window",
		"-t",
		sessionID,
		"-n",
		windowName,
		"-P",
		"-F",
		"\"#{window_id}--#{pane_id}\"",
	}
	stdout, stderr, err := t.executor.Execute(CommandName, args...)
	if err != nil {
		return nil, NewTmuxError(stdout, stderr, err)
	}
	output := stdout
	strings.TrimSpace(output)
	splitOutput := strings.Split(output, "--")
	if len(splitOutput) != 2 {
		log.Debug().Interface("output", splitOutput).Msg("invalid output from list-panes sub command")
		return nil, NewTmuxError(stdout, "", errors.New("cannot parse the windowID and pane ID from the list-panes output"))
	}
	return &TmuxCmdResponse{
		SessionID: "",
		WindowID:  splitOutput[0][1:],
		PaneID:    splitOutput[1][1:],
	}, nil
}

func (t *TmuxWrapper) newPane(targetPaneID string, sizeInPercentage int, horizontalSplit bool) (*TmuxCmdResponse, error) {
	// tmux splitw -h -p 10 -t 0 -P -F "#{pane_id}"
	// %10

	var args = []string{
		"splitw",
		"-h",
		"-p",
		strconv.Itoa(sizeInPercentage),
		"-t",
		targetPaneID,
		"-P",
		"-F",
		"\"#{window_id}--#{pane_id}\"",
	}
	if !horizontalSplit {
		args[1] = "-v"
	}
	stdout, stderr, err := t.executor.Execute(CommandName, args...)
	if err != nil {
		return nil, NewTmuxError(stdout, stderr, err)
	}
	output := stdout
	strings.TrimSpace(output)
	splitOutput := strings.Split(output, "--")
	if len(splitOutput) != 2 {
		log.Debug().Interface("output", splitOutput).Msg("invalid output from list-panes sub command")
		return nil, NewTmuxError(stdout, "", errors.New("cannot parse the windowID and pane ID from the list-panes output"))
	}
	return &TmuxCmdResponse{
		SessionID: "",
		WindowID:  splitOutput[0][1:],
		PaneID:    splitOutput[1][1:],
	}, nil
}

func (t TmuxWrapper) killSession(sessionID string) {
	// tmux kill-session -t session2
	log.Debug().Msgf("error while creating a new session, killing the session(%s) if it's created", sessionID)
	var args = []string{
		"kill-session",
		"-t",
		sessionID,
	}
	stdout, stderr, err := t.executor.Execute(CommandName, args...)
	if err != nil {
		log.Error().Err(err).Str("stdout", stdout).
			Str("stderr", stderr).
			Str("sessionID", sessionID).Msg("unable to kill the session")
	}
}

type ICommandExecutor interface {
	Execute(name string, args ...string) (string, string, error)
}

type CommandExecutor struct {
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (c *CommandExecutor) Execute(name string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.Command(name, args...)
	log.Debug().Str("command", command.String()).Msgf("executing...")
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err
}

type NOOPExecutor struct {
}

func NewNOOPExecutor() *NOOPExecutor {
	return &NOOPExecutor{}
}

func (c *NOOPExecutor) Execute(name string, args ...string) (string, string, error) {
	command := exec.Command(name, args...)
	log.Debug().Str("command", command.String()).Msgf("executing...")
	return "@5--%15", "", nil
}
