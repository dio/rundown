package runner

// Copyright 2022 Dhi Aurrahman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/dio/rundown/internal/archives"
	"github.com/tetratelabs/run"
)

// Run runs the prepared cmd, given an archive.
func Run(cmd *exec.Cmd, archive archives.Archive) (int, error) {
	err := cmd.Start()
	if err != nil {
		return 1, fmt.Errorf("failed to start %s: %w", archive.BinaryName(), err)
	}

	// Buffered, since caught by sigchanyzer: misuse of unbuffered os.Signal channel as argument to
	// signal.Notify.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		_ = cmd.Process.Signal(s)
		// TODO(dio): Handle windows.
	}()

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, _ := exitError.Sys().(syscall.WaitStatus)
			return waitStatus.ExitStatus(), nil
		}
		return 1, fmt.Errorf("failed to launch %s: %v", archive.BinaryName(), err)
	}
	return 0, nil
}

// MakeCmd returns initialized command.
func MakeCmd(binary string, args []string, out io.Writer) *exec.Cmd {
	cmd := exec.Command(binary, args...) //nolint:gosec
	cmd.Stdin = os.Stdin
	// TODO(dio): Kill all child processes.
	if out == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = out
	}
	cmd.Stderr = os.Stderr
	return cmd
}

// CanBeDisabled hold a service that can be disabled.
type CanBeDisabled struct {
	disabled bool
	g        *run.Group
	s        run.Service
}

// Manage delegate disabled checking to this object.
func (c *CanBeDisabled) Manage(g *run.Group, s run.Service, flags *run.FlagSet) {
	if g == nil {
		return
	}
	flags.BoolVar(
		&c.disabled,
		"disable-"+s.Name(),
		false,
		"Disable "+s.Name())
	c.g = g
	c.s = s
}

// IsTrue returns true when a managed service is disabled.
func (c *CanBeDisabled) IsTrue() bool {
	if c.g == nil || !c.disabled {
		return false
	}
	c.g.Deregister(c.s)
	return true
}
