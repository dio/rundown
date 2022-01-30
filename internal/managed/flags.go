package managed

import (
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tetratelabs/run"
)

// Managed holds common flags that can be shared across services.
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

	// --<name>-config. For example: --proxy-config.
	flags.StringVar(
		&m.ConfigFile,
		s.Name()+"-config",
		m.ConfigFile,
		"Path to the "+title+" config file")

	// --<name>-version. For example: --proxy-version.
	flags.StringVar(
		&m.Version,
		s.Name()+"-version",
		m.DefaultVersion,
		title+" version",
	)

	// --<name>-directory. For example: --proxy-directory.
	flags.StringVar(
		&m.Dir,
		s.Name()+"-directory",
		os.Getenv(strcase.ToScreamingSnake(s.Name())+"_HOME"),
		"Path to the "+title+" work directory",
	)

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
