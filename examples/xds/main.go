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

package main

import (
	_ "embed" // to allow embedding files.
	"fmt"
	"os"

	bootstrapv3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3" // allow to resolve envoy.extensions.upstreams.http.v3.*
	"github.com/tetratelabs/run"
	runsignal "github.com/tetratelabs/run/pkg/signal"
	"github.com/tetratelabs/telemetry"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"

	"github.com/dio/rundown/api/proxy"
	"github.com/dio/rundown/api/xds"
	configv1 "github.com/dio/rundown/generated/xds/config/v1"
)

//go:embed proxy.yaml
var configYAML []byte

func main() {
	var (
		logger            = telemetry.NoopLogger()
		g                 = &run.Group{Name: "example", Logger: logger}
		xdsServerConfig   = &configv1.Config{Host: "localhost", Port: 8080}
		xdsServer         = xds.New(g, &xds.Config{Logger: g.Logger, Config: xdsServerConfig})
		proxyServerConfig = &proxy.Config{Logger: g.Logger, GenerateConfig: generate}
		proxyServer       = proxy.New(g, proxyServerConfig)
		signalHandler     = new(runsignal.Handler)
	)
	g.Register(xdsServer, proxyServer, signalHandler)
	if err := g.Run(); err != nil {
		fmt.Printf("program exit: %+v\n", err)
		os.Exit(1)
	}
}

func generate() (*bootstrapv3.Bootstrap, error) {
	j, err := yaml.YAMLToJSON(configYAML)
	if err != nil {
		return nil, err
	}
	var config bootstrapv3.Bootstrap
	err = protojson.Unmarshal(j, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
