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

package xds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryservice "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	ratelimitrunner "github.com/envoyproxy/ratelimit/src/service_cmd/runner"
	"github.com/tetratelabs/run"
	"github.com/tetratelabs/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"

	configv1 "github.com/dio/rundown/generated/xds/config/v1"
	"github.com/dio/rundown/internal/managed"
	"github.com/dio/rundown/internal/xds"
)

// Config holds the configuration object for running the rate-limit service.
type Config struct {
	Logger telemetry.Logger
	Config *configv1.Config
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
			Titleize: func(string) string {
				return "xDS Service"
			},
		},
	}
}

// Service is a run.Service implementation that runs the rate-limit service.
type Service struct {
	cfg        *Config
	g          *run.Group
	runner     *ratelimitrunner.Runner
	managed    *managed.Flags
	server     server.Server
	grpcServer *grpc.Server
	cache      cache.SnapshotCache
}

var _ run.Config = (*Service)(nil)

// Name returns the service name.
func (s *Service) Name() string {
	return "xds-service"
}

// FlagSet provides command line flags for the service.
func (s *Service) FlagSet() *run.FlagSet {
	flags := run.NewFlagSet("xDS Service options")
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

		var cfg configv1.Config
		if err = protojson.Unmarshal(b, &cfg); err != nil {
			return err
		}
		s.cfg.Config = &cfg
	}

	if s.cfg.Config == nil {
		return errors.New("xds service config is required")
	}
	return s.cfg.Config.ValidateAll()
}

// PreRun prepares the service to run.
func (s *Service) PreRun() (err error) {
	if s.managed.IsDisabled() {
		return nil
	}

	ctx := context.Background()

	// TODO(dio): Change this implementation.
	s.cache = cache.NewSnapshotCache(false, cache.IDHash{}, logger{})
	snapshot := xds.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		return err
	}

	if err := s.cache.SetSnapshot(ctx, "test-id", snapshot); err != nil {
		return err
	}

	s.server = server.NewServer(ctx, s.cache, &callbacks{})

	// TODO(dio): Take this grpc server setup to a dedicated module.
	s.grpcServer = grpc.NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Config.Host, s.cfg.Config.Port))
	if err != nil {
		return err
	}
	// This probably requires refactoring. It depends on how we do tests.
	return s.registerAndServe(listener)
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
	s.grpcServer.GracefulStop()
}

func (s *Service) registerAndServe(listener net.Listener) error {
	// Register services.
	discoveryservice.RegisterAggregatedDiscoveryServiceServer(s.grpcServer, s.server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(s.grpcServer, s.server)
	clusterservice.RegisterClusterDiscoveryServiceServer(s.grpcServer, s.server)
	routeservice.RegisterRouteDiscoveryServiceServer(s.grpcServer, s.server)
	listenerservice.RegisterListenerDiscoveryServiceServer(s.grpcServer, s.server)
	secretservice.RegisterSecretDiscoveryServiceServer(s.grpcServer, s.server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(s.grpcServer, s.server)
	return s.grpcServer.Serve(listener)
}

type logger struct{}

func (l logger) Debugf(format string, args ...interface{}) {}
func (l logger) Infof(format string, args ...interface{})  {}
func (l logger) Warnf(format string, args ...interface{})  {}
func (l logger) Errorf(format string, args ...interface{}) {}

type callbacks struct{}

func (c *callbacks) Report() {}
func (c *callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	return nil
}
func (c *callbacks) OnStreamClosed(id int64) {}
func (c *callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	return nil
}
func (c *callbacks) OnDeltaStreamClosed(id int64) {}
func (c *callbacks) OnStreamRequest(int64, *discoveryservice.DiscoveryRequest) error {
	return nil
}
func (c *callbacks) OnStreamResponse(context.Context, int64, *discoveryservice.DiscoveryRequest, *discoveryservice.DiscoveryResponse) {
}
func (c *callbacks) OnStreamDeltaResponse(id int64, req *discoveryservice.DeltaDiscoveryRequest, res *discoveryservice.DeltaDiscoveryResponse) {
}
func (c *callbacks) OnStreamDeltaRequest(id int64, req *discoveryservice.DeltaDiscoveryRequest) error {
	return nil
}
func (c *callbacks) OnFetchRequest(_ context.Context, req *discoveryservice.DiscoveryRequest) error {
	return nil
}
func (c *callbacks) OnFetchResponse(*discoveryservice.DiscoveryRequest, *discoveryservice.DiscoveryResponse) {
}
