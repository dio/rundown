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

package managed

import (
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tetratelabs/run"
)

// Flags holds common flags that can be shared across services.
type Flags struct {
	DefaultVersion string
	Version        string
	Dir            string
	ConfigFile     string

	disabled bool
	g        *run.Group
	s        run.Service
}

// Manage delegate disabled checking to this object.
func (m *Flags) Manage(flags *run.FlagSet, g *run.Group, s run.Service) {
	if g == nil {
		return
	}
	title := titleize(s.Name())

	// When default version is not defined, we should not register the version flag, since it is
	// probably hardcoded or not relevant (e.g. for internal implementation).
	if m.DefaultVersion != "" {
		// --<name>-version. For example: --proxy-version.
		flags.StringVar(
			&m.Version,
			s.Name()+"-version",
			m.DefaultVersion,
			title+" version",
		)
	}

	// --<name>-directory. For example: --proxy-directory.
	flags.StringVar(
		&m.Dir,
		s.Name()+"-directory",
		os.Getenv(strcase.ToScreamingSnake(s.Name())+"_HOME"),
		"Path to the "+title+" work directory",
	)

	// --<name>-config. For example: --proxy-config.
	flags.StringVar(
		&m.ConfigFile,
		s.Name()+"-config",
		m.ConfigFile,
		"Path to the "+title+" config file")

	// --disable-<name>. For example: --disable-proxy.
	flags.BoolVar(
		&m.disabled,
		"disable-"+s.Name(),
		false,
		"Disable "+title)

	m.g = g
	m.s = s
}

// IsDisabled returns true when a managed service is disabled.
func (m *Flags) IsDisabled() bool {
	if m.g == nil || !m.disabled {
		return false
	}
	m.g.Deregister(m.s)
	return true
}

// titleize properly capitalize kebab case to title case.
func titleize(name string) string {
	return strings.Title(strings.Join(strings.Split(name, "-"), " "))
}
