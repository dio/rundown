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

	"github.com/tetratelabs/run"
	runsignal "github.com/tetratelabs/run/pkg/signal"
	"github.com/tetratelabs/telemetry"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"

	"github.com/dio/rundown/api/ratelimit"
	settingsv1 "github.com/dio/rundown/generated/ratelimit/settings/v1"
)

//go:embed ratelimit.yaml
var configYAML []byte

func main() {
	var (
		logger    = telemetry.NoopLogger()
		g         = &run.Group{Name: "example", Logger: logger}
		rateLimit = ratelimit.New(g, &ratelimit.Config{
			Logger:         g.Logger,
			GenerateConfig: generate,
		})
		signalHandler = new(runsignal.Handler)
	)
	g.Register(rateLimit, signalHandler)
	if err := g.Run(); err != nil {
		fmt.Printf("program exit: %+v\n", err)
		os.Exit(1)
	}
}

func generate() (*settingsv1.Settings, error) {
	j, err := yaml.YAMLToJSON(configYAML)
	if err != nil {
		return nil, err
	}
	var config settingsv1.Settings
	err = protojson.Unmarshal(j, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
