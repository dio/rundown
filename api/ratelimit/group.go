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

package ratelimit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3" // added to resolve v3.HttpConnectionManager.
	ratelimitrunner "github.com/envoyproxy/ratelimit/src/service_cmd/runner"
	"github.com/tetratelabs/run"
	"github.com/tetratelabs/telemetry"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"

	"github.com/dio/rundown/generated/ratelimit/settings"
	"github.com/dio/rundown/internal/managed"
	"github.com/dio/rundown/internal/ratelimit"
)

// Config holds the configuration object for running the rate-limit service.
type Config struct {
	Logger   telemetry.Logger
	Settings *settings.Settings
}

// New returns a new run.Service that wraps rate-limit service. Setting the cfg to nil, expecting
// setting the rate-limit settings object from a file.
func New(g *run.Group, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{} // TODO(dio): Have a way to generate default config.
	}
	return &Service{
		cfg: cfg,
		g:   g,
		managed: &managed.Flags{
			CanBeDisabledOnly: true, // Allow for disable flag generation and checking only.
		},
	}
}

// Service is a run.Service implementation that runs the rate-limit service.
type Service struct {
	cfg     *Config
	g       *run.Group
	runner  *ratelimitrunner.Runner
	managed *managed.Flags
}

var _ run.Config = (*Service)(nil)

// Name returns the service name.
func (s *Service) Name() string {
	return "rate-limit-service"
}

// FlagSet provides command line flags for the service.
func (s *Service) FlagSet() *run.FlagSet {
	flags := run.NewFlagSet("Rate Limit Service options")
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

		// Probably a .yaml file. We simply check the extension here.
		if filepath.Ext(s.managed.ConfigFile) == ".yaml" || filepath.Ext(s.managed.ConfigFile) == ".yml" {
			b, err = yaml.YAMLToJSON(b)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
		}

		var cfg settings.Settings
		if err = protojson.Unmarshal(b, &cfg); err != nil {
			return err
		}
		s.cfg.Settings = &cfg
	}

	if s.cfg.Settings == nil {
		return errors.New("proxy config is required")
	}
	return s.cfg.Settings.ValidateAll()
}

// PreRun prepares the service to run.
func (s *Service) PreRun() error {
	if s.managed.IsDisabled() {
		return nil
	}
	runner := ratelimitrunner.NewRunner(ratelimit.NewSettings(s.cfg.Settings))
	s.runner = &runner
	return nil
}

// Serve runs the service.
func (s *Service) Serve() (err error) {
	defer func() {
		recovered := recover()
		if err != nil {
			err = recovered.(error)
		}
	}()
	s.runner.Run()
	return
}

// GracefulStop stops the underlying process by sending interrupt.
func (s *Service) GracefulStop() {
	if s.runner != nil {
		s.runner.Stop()
	}
}
