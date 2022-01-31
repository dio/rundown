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

package ratelimit_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sigs.k8s.io/yaml"

	"github.com/dio/rundown/generated/ratelimit/settings"
	"github.com/dio/rundown/internal/ratelimit"
)

//go:embed testdata/settings.yaml
var config []byte

func TestNewSettings(t *testing.T) {
	s := &settings.Settings{
		Host: &wrapperspb.StringValue{Value: "127.0.0.1"},
	}
	c := ratelimit.NewSettings(s)
	require.Equal(t, c.Host, s.Host.Value)

	j, err := yaml.YAMLToJSON(config)
	require.NoError(t, err)
	require.NotNil(t, j) // Because when "j" is nil, it is not an error.

	var s1 settings.Settings
	require.NoError(t, protojson.Unmarshal(j, &s1))

	c1 := ratelimit.NewSettings(&s1)
	require.Equal(t, c1.Host, s1.Host.Value)
	require.Equal(t, c1.Port, int(s1.Port.Value))
	require.Equal(t, c1.GrpcHost, s1.GrpcHost.Value)
	require.Equal(t, c1.GrpcPort, int(s1.GrpcPort.Value))
}
