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

package extauthz

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/tetratelabs/run"
	"github.com/tetratelabs/telemetry"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/dio/rundown/generated/authservice/config"
	"github.com/dio/rundown/internal/archives"
	"github.com/dio/rundown/internal/downloader"
	"github.com/dio/rundown/internal/managed"
	"github.com/dio/rundown/internal/runner"
)

var (
	// Default binary version.
	DefaultBinaryVersion = "0.6.0-rc0"
	// Default download timeout.
	DefaultDownloadTimeout = 30 * time.Second
)

// Config holds the configuration object for running the auth_server.
type Config struct {
	Logger       telemetry.Logger
	FilterConfig *config.Config
}

// New returns a new run.Service that wraps auth_server binary. Setting the cfg to nil, expecting
// setting the auth_server's --filter_config from a file.
func New(g *run.Group, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{} // TODO(dio): Have a way to generate default config.
	}
	return &Service{
		cfg:     cfg,
		g:       g,
		archive: &archives.ExtAuthz{},
		managed: &managed.Flags{
			DefaultVersion: DefaultBinaryVersion,
		},
	}
}

// Service is a run.Service implementation that runs auth_server.
type Service struct {
	cfg     *Config
	g       *run.Group
	archive *archives.ExtAuthz
	managed *managed.Flags
	cmd     *exec.Cmd
}

var _ run.Config = (*Service)(nil)

// Name returns the service name.
func (s *Service) Name() string {
	return "external-auth-service"
}

// FlagSet provides command line flags for external auth-service.
func (s *Service) FlagSet() *run.FlagSet {
	flags := run.NewFlagSet("External AuthN/AuthZ Service options")
	s.managed.Manage(flags, s.g, s)
	return flags
}

// Validate validates the given configuration.
func (s *Service) Validate() error {
	if s.managed.IsDisabled() {
		return nil
	}

	if s.managed.ConfigFile != "" {
		b, err := os.ReadFile(s.managed.ConfigFile)
		if err != nil {
			return err
		}
		var cfg config.Config
		if err = protojson.Unmarshal(b, &cfg); err != nil {
			return err
		}
		s.cfg.FilterConfig = &cfg
	}

	if s.cfg.FilterConfig == nil {
		return errors.New("filter config is required")
	}
	return s.cfg.FilterConfig.ValidateAll()
}

// PreRun prepares the binary to run.
func (s *Service) PreRun() (err error) {
	if s.managed.Dir == "" {
		// To make sure we have a work directory.
		dir, err := ioutil.TempDir("", s.archive.BinaryName())
		if err != nil {
			return nil
		}
		s.managed.Dir = dir
	}

	if s.managed.Version != "" {
		s.archive.VersionUsed = s.managed.Version
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultDownloadTimeout)
	defer cancel()

	// Check and download the versioned binary.
	binaryPath, err := downloader.DownloadVersionedBinary(ctx, s.archive, s.managed.Dir)
	if err != nil {
		return err
	}

	// Generate JSON config to run the auth_server. See: authservice/docs/README.md.
	jsonConfig, err := protojson.Marshal(s.cfg.FilterConfig)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(s.managed.Dir, "*.json")
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if _, err = tmp.Write(jsonConfig); err != nil {
		return err
	}

	configPath := tmp.Name() // effective config path.
	s.cmd = runner.MakeCmd(binaryPath, []string{"--filter_config", configPath}, os.Stdout)
	return nil
}

// Serve runs the binary.
func (s *Service) Serve() error {
	// Run the downloaded auth_server with the generated config in s.configPath.
	if exitCode, err := runner.Run(s.cmd, s.archive); err != nil {
		s.cfg.Logger.Error(fmt.Sprintf("%s exit with %d", s.archive.BinaryName(), exitCode), err)
		return err
	}
	return nil
}

// GracefulStop stops the underlying process by sending interrupt.
func (s *Service) GracefulStop() {
	if s.cmd != nil {
		s.cmd.Process.Signal(os.Interrupt)
	}
}
